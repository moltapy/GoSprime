import os, re, sys
import signal, psutil
import argparse, platform
import yaml, yaml.scanner, subprocess

from pathlib import Path
from loguru import logger
from typing import Dict, Tuple, List
from concurrent.futures import ProcessPoolExecutor


KEYS        =   ['apps', 'paths', 'names', 'ranges', 'formats', 'tags', 'parallel']
APPS        =   ['ExtractSamples', 'ExtractAutos', 'ConcatAutos', 'GenerateScores', 'MappingArchs']
PATHS       =   ['bcftool_filepath', 'sprimejar_filepath', 'geneticmap_filepath', 'maskbed_filepath', 
                 'origin_samplelist_filepath', 'splited_samplelists_dirpath', 'outgroup_samplelists_dirpath', 
                 'modernhuman_vcffiles_filepath', 'modernhuman_extracted_vcffiles_dirpath', 'concated_vcffiles_dirpath', 
                 'concated_genotypes_filepath', 'generated_scorefiles_dirpath', 'archaic_vcffile_x_filepath', 
                 'archaic_vcffile_y_filepath', 'summary_script_filepath', 'contour_script_filepath']
NAMES       =   ['outgroup_name', 'modernhuman_vcffiles_name', 'column_name_x', 'column_name_y']
RANGES      =   ['extracted_samples_range', 'extracted_autos_chroms_range', 'extracted_autos_groups_range', 
                 'concat_autos_groups_range', 'generate_scores_groups_range', 'generate_scores_chroms_range']
FORMATS     =   ['samplelist_lines_separators']
TAGS        =   ['samplelists_isheader', 'extracted_vcfs_isoverwrite', 'concat_vcfs_isoverwrite', 
                 'generate_scores_isoverwrite', 'maskbed_isexclude', 'depth_x', 'depth_y']
PARALLEL    =   ['cores', 'threads']


def register_logger(debug: bool) -> None:

    global logger
    level = "DEBUG" if debug else "INFO"
    isTerminal = sys.stdout.isatty()
    logger.remove()
    logger.add(
        sink = sys.stdout,
        format = (
            "<green>{time:YYYY-MM-DD HH:mm:ss.SSS}</green> | "
            "<level>{level: <8}</level> | "
            "<cyan>{extra[module]}</cyan> - "
            "<white>{message}</white>"
        ),
        colorize = isTerminal,
        level = level
    )
    logger   = logger.bind(module = "SPRIME_MAIN")


executor = None
current_process = psutil.Process()


def parse_args():

    parser = argparse.ArgumentParser(description = "Main App of Gosprime")
    parser.add_argument("-c", "--config", type = str, default = "config.yaml", help = "Path of config files")
    parser.add_argument("-w", "--workdir", type = str, default = Path(__file__).parent, help = "Directory of executable dirs")
    parser.add_argument("-d", "--debug", action= "store_true", default = False, help = "Debug mode or not")
    args = parser.parse_args()
    return args


def check_configs(configs: Dict[str,str]):
    
    global logger
    logger = logger.bind(module = "CONFIG_CHECK")
    logger.info("Starting configuration validation check ... ")

    for key in KEYS:
        if key not in configs.keys():
            logger.error(f"Key: {key} is missing, please check!")
            exit(-1)

    for app in APPS:
        if app not in configs["apps"]:
            logger.error(f"Key: apps.{app} is missing, please check!")
            exit(-1)

    for path in PATHS:
        if path not in configs["paths"].keys():
            logger.error(f"Key: paths.{path} is missing, please check!")
            exit(-1)

    for name in NAMES:
        if name not in configs["names"].keys():
            logger.error(f"Key: names.{key} is missing, please check!")
            exit(-1)

    for rang in RANGES:
        if rang not in configs["ranges"].keys():
            logger.error(f"Key: ranges.{key} is missing, please check!")
            exit(-1)

    for form in FORMATS:
        if form not in configs["formats"].keys():
            logger.error(f"Key: formats.{key} is missing, please check!")
            exit(-1)

    for tag in TAGS:
        if tag not in configs["tags"].keys():
            logger.error(f"Key: tags.{key} is missing, please check!")
            exit(-1)

    for parallel in PARALLEL:
        if parallel not in configs["parallel"].keys():
            logger.error(f"Key: parallel.{key} is missing, please check!")
            exit(-1)

    logger.info("Verification completed, proceeding to next stage ... ")
    

def load_configs(conf: str) -> Tuple[str]:

    global logger
    logger = logger.bind(module = "CONFIG_LOAD")
    logger.info(f"Start configuration loading process, reading from {conf} ... ")

    try:
        configs = yaml.full_load(open(conf,"rb"))
    except yaml.scanner.ScannerError as e:
        logger.error(f"Problem occurred when loading configs, reason: {e}, please check!")
        exit(-1)

    check_configs(configs) if configs else None

    logger.info("Loading completed, proceeding to next stage ... ")
    return (configs["apps"],configs["paths"],
            configs["names"],configs["ranges"],
            configs["formats"],configs["tags"],configs["parallel"])


def set_environs(binpath: str, apps: List[str]):

    global logger
    logger = logger.bind(module = "ENVIRON_SET")
    logger.info("Start executables building process ... ")

    filenames = os.scandir(binpath)
    dirnames  = [file.name for file in filenames if file.is_dir()]
    system    = platform.system()
    for app in apps:
        if app not in dirnames:
            logger.error(f"Executable <{app}> not in working directory {binpath}, please check!")
            exit(-1)
        else:
            os.chdir(Path(binpath) / app)
            appname = app + ".exe" if system.lower() == "windows" else app
            approute = Path(binpath) / app / appname
            if Path.exists(approute):
                logger.info(f"Executable <{appname}> already exists in {Path(binpath) / app}, continue ... ")
            else:
                logger.info(f"Executable <{appname}> not exists in {Path(binpath) / app}, building ... ")
                status = subprocess.run("go build", shell = True, capture_output = True, text = True)
                if status.returncode == 0:
                    logger.info(f"Executable <{appname}> compiled without errors, continue...")
                else:
                    logger.info(f"Executable <{appname}> compiled with errors, reason: {status.stderr}")

    os.chdir(binpath)
    logger.info("Environment set completed, proceeding to next stage ... ")


def init_worker():
    signal.signal(signal.SIGINT, signal.SIG_IGN)


def extract_samples(binpath: str, binname: str, samplelist: str, otgname: str, splitedpath: str,
                    otgpath: str, limits: List[str] = None, title: bool = True,fstr: str = None):
    
    global logger
    logger   = logger.bind(module = "EXTRACT_SAMPLES")
    
    appname  = f"{binname}.exe" if platform.system().lower() == "windows" else binname
    approute = Path(binpath) / binname / appname
    
    command  = f"{approute} -s {samplelist} -g {otgname} -p {splitedpath} -o {otgpath}"
    command += f" -l {limits}"  if limits else ""    
    command +=  " -t"           if title else ""
    command += f" -f {fstr}"    if fstr else "" 

    logger.info(f"Executable <{appname}> starts processing, outgroup : {otgname}")
    logger.debug(f"Command: {command}")
    status = subprocess.run(command, shell = True, text = True, stderr = subprocess.STDOUT)
    if status.returncode == 0:
        logger.info(f"Executable <{appname}> processing completed, proceeding to next stage ... ")
    else:
        logger.error(f"Executable <{appname}> processing wrong, terminated, see above for details")
        exit(-1)


def extract_autos(binpath: str, binname: str, samplepath: str, modernpath: str,
                  modernout: str, vcfname: str = None, overwrite:bool = False,
                  autosrange: List[int] = None, groupsrange: List[str] = None, 
                  bcftoolpath: str = None, threads: int = None, cores: int = None):
    
    global logger
    logger   = logger.bind(module = "EXTRACT_AUTOS")

    appname  = f"{binname}.exe" if platform.system().lower() == "windows" else binname
    approute = Path(binpath) / binname / appname

    command  = f"{approute} -s {samplepath} -m {modernpath} -o {modernout}"
    command += f" -b {bcftoolpath}" if bcftoolpath else ""
    command += f" -n {vcfname}"     if vcfname else ""
    command += f" -a {autosrange}"  if autosrange else ""
    command += f" -g {groupsrange}" if groupsrange else ""
    command += f" -t {threads}"     if threads else ""
    command += f" -c {cores}"       if cores else ""
    command +=  " --overwrite"               if overwrite else ""

    logger.info(f"Executable <{appname}> starts processing ... ")
    logger.debug(f"Command: {command}")
    status = subprocess.run(command, shell = True, text = True, stderr = subprocess.STDOUT)
    if status.returncode == 0:
        logger.info(f"Executable <{appname}> processing completed, proceeding to next stage ... ")
    else:
        logger.error(f"Executable <{appname}> processing wrong, terminated, see above for details")
        exit(-1)


def concat_autos(binpath: str, binname: str, fagrpspath: str, 
                 concatedpath: str, bcftoolpath: str = None, 
                 groupsrange: List[str] = None, overwrite: bool = False):
    
    global logger
    logger   = logger.bind(module = "CONCAT_AUTOS")

    appname  = f"{binname}.exe" if platform.system().lower() == "windows" else binname
    approute = Path(binpath) / binname / appname

    command  = f"{approute} -d {fagrpspath} -o {concatedpath}"
    command += f" -g {groupsrange}" if groupsrange else ""
    command += f" -b {bcftoolpath}" if bcftoolpath else ""
    command +=  " --overwrite"               if overwrite else ""

    logger.info(f"Executable <{appname}> starts processing ... ")
    logger.debug(f"Command: {command}")
    status = subprocess.run(command, shell = True, text = True, stderr = subprocess.STDOUT)
    if status.returncode == 0:
        logger.info(f"Executable <{appname}> processing completed, proceeding to next stage ... ")
    else:
        logger.error(f"Executable <{appname}> processing wrong, terminated, see above for details")
        exit(-1)


def generate_scores(binpath: str, binname: str, genopath: str, 
                    maproute: str, jarpath: str, outgrouplist: str, 
                    otgrpname: str, scorefpath: str, groupsrange: List[str] = None, 
                    chromsrange: List[int] = None, threads: int = None, cores: int = None, overwrite: bool = False):
    
    global logger
    logger   = logger.bind(module = "GEN_SCORES")

    appname  = f"{binname}.exe" if platform.system().lower() == "windows" else binname
    approute = Path(binpath) / binname / appname

    command  = f"{approute} -g {genopath} -m {maproute} -j {jarpath} -u {outgrouplist} -n {otgrpname} -o {scorefpath}"
    command += f" -r {groupsrange}" if groupsrange else ""
    command += f" -l {chromsrange}" if chromsrange else ""
    command += f" -t {threads}"     if threads else ""
    command += f" -c {cores}"       if cores else ""
    command +=  " --overwrite"               if overwrite else ""

    logger.info(f"Executable <{appname}> starts processing ... ")
    logger.debug(f"Command: {command}")
    status = subprocess.run(command, shell = True,text = True, stderr = subprocess.STDOUT)
    if status.returncode == 0:
        logger.info(f"Executable <{appname}> processing completed, proceeding to next stage ... ")
    else:
        logger.error(f"Executable <{appname}> processing wrong, terminated, see above for details")
        exit(-1)


def mapping_archs(binpath: str, binname: str, maskpath: str, 
                  archpath: str, scorepath: str, logfile: str,
                  arrayname: str, sep: str = None, reverse: bool = False, depth: bool = False):
    
    global logger
    logger   = logger.bind(module = "MAPPING_ARCH")

    appname  = f"{binname}.exe" if platform.system().lower() == "windows" else binname
    approute = Path(binpath) / binname / appname

    command  = f"{approute} -m {maskpath} -v {archpath} -p {scorepath} -n {arrayname}"
    command += f" -s {sep}" if sep else ""
    command += " --reverse" if reverse else ""
    command += " --depth"   if depth else ""

    logger.debug(f"Command: {command}, Log: {logfile}")
    status = subprocess.run(command, stderr = open(logfile, "at", encoding = "utf-8"), shell = True)
    if status.returncode != 0:
        logger.error(f"Fail to process {scorepath}, column: {arrayname}, reason: {status.stderr}")
        logger.error(f"Error logs in log: {logfile}")
        exit(-1)
    logger.debug(f"Command: {command} processing completed, log: {logfile}")


def check_path(placeholder: str,path: str) -> int:
    
    global logger
    logger.bind(module = "CHECK_PATH")
    logger.info(f"Starting file path validation check of {path} ... ")
    return path.index("{"+ placeholder + "}")


def process_chromosome(workmem: str, appname: str, paths: Dict[str,str], 
                       names: Dict[str,str], tags: Dict[str,str], file: str, chrom: str):

    global logger
    logger = logger.bind(module = "PROCESS_CHROM")

    try:       
        maskfile    = paths["maskbed_filepath"].replace("{chr}", str(chrom))
        archaicvcfx = paths["archaic_vcffile_x_filepath"].replace("{chr}", str(chrom))
        archaicvcfy = paths["archaic_vcffile_y_filepath"].replace("{chr}", str(chrom))
        logfile     = Path.joinpath(file.parent,"".join(file.name.split(".")[:len(file.name.split("."))-1])+"_match.log")
        open(logfile,"w") 

        mapping_archs(
            workmem, 
            appname,
            maskfile,
            archaicvcfx,
            file,
            logfile,
            names["column_name_x"],
            r"\\t",
            tags["maskbed_isexclude"],
            tags["depth_x"]
        )
        
        mapping_archs(
            workmem, 
            appname,
            maskfile,
            archaicvcfy,
            file,
            logfile,
            names["column_name_y"],
            r"\\t",
            tags["maskbed_isexclude"],
            tags["depth_y"]
        )
        
        logger.info(f"Complete mapping chromosome {chrom}, file changed: {file}, log: {logfile}")
    except Exception as e:
        logger.error(f"Chromosome {chrom} processing failed, reason: {str(e)}")
        raise  


def contour_plot(paths: Dict[str,str], names: Dict[str,str]) -> None:

    global logger
    logger = logger.bind(module = "CONTOUR_PLOT")

    groups = os.scandir(paths["generated_scorefiles_dirpath"])
    for group in groups:
        if group.is_dir():
            logger.info(f"Start processing {group.name} ... ")
            summary = Path.joinpath(Path(group),"summary.txt")
            logger.info(f"Start summarizing {group.name}, result: {str(summary)}")
            command = f"Rscript {paths['summary_script_filepath']} {group.path} {str(summary)}"
            logger.debug(f"Command: {command}")
            status = subprocess.run(command, shell = True, capture_output = True, text = True)
            if status.returncode == 0:
                logger.info(f"Complete summarizing {group.name}")
            else:
                logger.error(f"Summarizing {group.name} error, reason: {status.stderr}")
                exit(-1)

            plot = Path.joinpath(Path(group),f"{group.name}_contour")
            logger.info(f"Start plotting {group.name}, result: {str(plot)}.png")
            command = f"Rscript {paths['contour_script_filepath']} {str(summary)} {str(plot)} {names['column_name_x']} {names['column_name_y']}"
            logger.debug(f"Command: {command}")
            status = subprocess.run(command, shell = True, capture_output= True, text = True)
            if status.returncode == 0:
                logger.info(f"Complete plotting {group.name}")
            else:
                logger.error(f"Plotting {group.name} error, reason: {status.stderr}")
                os.exit(-1)
            logger.info(f"Complete processing {group.name}")


def main():

    global logger
    args = parse_args()
    register_logger(args.debug)
    logger.info(f"Directory: {args.workdir}, config: {args.config}")
    apps, paths,names,ranges,formats,tags,parallel = load_configs(args.config)
    set_environs(args.workdir,apps)

    extract_samples(args.workdir, apps[0], paths["origin_samplelist_filepath"], names["outgroup_name"], 
                    paths["splited_samplelists_dirpath"], paths["outgroup_samplelists_dirpath"], ranges["extracted_samples_range"], 
                    tags["samplelists_isheader"], formats["samplelist_lines_separators"])
    
    extract_autos(args.workdir, apps[1], paths["splited_samplelists_dirpath"], paths["modernhuman_vcffiles_filepath"], 
                  paths["modernhuman_extracted_vcffiles_dirpath"], names["modernhuman_vcffiles_name"], 
                  tags["extracted_vcfs_isoverwrite"], ranges["extracted_autos_chroms_range"], 
                  ranges["extracted_autos_groups_range"], paths["bcftool_filepath"], parallel["threads"], parallel["cores"])

    concat_autos(args.workdir, apps[2], paths["modernhuman_extracted_vcffiles_dirpath"], paths["concated_vcffiles_dirpath"], 
                paths["bcftool_filepath"], ranges["concat_autos_groups_range"], tags["concat_vcfs_isoverwrite"])

    generate_scores(args.workdir, apps[3], paths["concated_genotypes_filepath"], paths["geneticmap_filepath"], 
                    paths["sprimejar_filepath"], f"{paths['outgroup_samplelists_dirpath']}/outgroup.txt", 
                    names["outgroup_name"], paths["generated_scorefiles_dirpath"], ranges["generate_scores_groups_range"],
                    ranges["generate_scores_chroms_range"], parallel["threads"], parallel["cores"], tags["generate_scores_isoverwrite"])

    try:
        logger = logger.bind(module = "SPRIME_MAIN")
        logger.info("Preparing for mapping, checking path ... ")
        check_path("chr",paths["maskbed_filepath"])
        check_path("chr",paths["archaic_vcffile_x_filepath"])
    except ValueError as e :
        logger.error(f"Invalid path detected, reason: {e}")
        exit(-1)

    pattern = re.compile(r"(\D*)(\d+)(\D*)")

    try:
        for fileobj in os.scandir(paths["generated_scorefiles_dirpath"]):
            files = list(os.scandir(Path(paths["generated_scorefiles_dirpath"]) / fileobj.name))
            scorelist = []
            if fileobj.is_dir() and list(os.scandir(Path(paths["generated_scorefiles_dirpath"]) / fileobj.name)):
                logger.info(f"Group {fileobj.name} detected, proceeding to next stage ... ")
                for file in files:
                    if file.name.endswith("score"):
                        scorelist.append(file.name)
                scorelist = sorted(scorelist,key = lambda x : int(pattern.findall(x)[0][1]))
                filelist = [Path(paths["generated_scorefiles_dirpath"]) / fileobj.name / x for x in scorelist]

                global executor
                executor = ProcessPoolExecutor(max_workers=22, initializer=init_worker)

                try:
                    futures = []
                    for index, file in enumerate(filelist):
                        chrom = index + 1
                        future = executor.submit(
                            process_chromosome,
                            args.workdir,
                            apps[4],
                            paths,
                            names,
                            tags,
                            file,
                            chrom
                        )
                        futures.append(future)
                    logger.info(f"Group {fileobj.name} submitted {len(futures)} tasks, columns: {names['column_name_x']} - {names['column_name_y']}, continue ... ")

                    for future in futures:

                        try:
                            future.result()

                        except Exception as e:
                            logger.error(f"Problem occurred when processing group: {fileobj.name}, reason: {str(e)}")

                except KeyboardInterrupt:

                    logger.error("Received KeyboardInterrupt, terminating all processes...")
                    children = current_process.children(recursive=True)

                    for child in children:

                        try:
                            child.kill()
                        
                        except psutil.NoSuchProcess:
                            pass
                    
                    os._exit(1)

                except Exception as e:
                    logger.error(f"Problem occurred when terminating, reason: {e}")

                finally:
                    executor.shutdown(wait = False)

            else:
                logger.info(f"Item {fileobj.name} is not a directory or empty, continue ... ")

    except KeyboardInterrupt:
        logger.error("Main process interrupted, cleaning up ... ")

    finally:

        children = current_process.children(recursive=True)
        
        for child in children:
        
            try:
                child.kill()
        
            except psutil.NoSuchProcess:
                pass
        
        if executor:
            executor.shutdown()
    
    logger.bind(module = "SPRIME_MAIN")
    logger.info("Start summarizing and plotting ...")
    contour_plot(paths, names)
    logger.info(f"Complete summarizing and plotting processes, contours in {str(paths['generated_scorefiles_dirpath'])}")

if __name__ == "__main__":
    main()