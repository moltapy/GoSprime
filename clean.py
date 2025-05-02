import yaml
import shutil
import argparse


def clean_total():

    paths = yaml.load(open("config.yaml", "rb"), yaml.FullLoader)["paths"]
    shutil.rmtree(paths["splited_samplelists_dirpath"])
    shutil.rmtree(paths["outgroup_samplelists_dirpath"])
    shutil.rmtree(paths["modernhuman_extracted_vcffiles_dirpath"])
    shutil.rmtree(paths["concated_vcffiles_dirpath"])
    shutil.rmtree(paths["generated_scorefiles_dirpath"])
    
    
def clean_mapping():

    for index in range(1, 23):
        path = f"score_groups/CDX/CDX_YRI_{index}.score"
        movepath = f"score_groups/CDX/CDX_YRI_{index}_clean.score"
        with open(path, "rt") as infile:
            outfile = open(movepath,"wt")
            for line in infile:
                line = line.strip().split("\t")
                line = "\t".join(line[:8])
                outfile.write(line+"\n")
            outfile.close()
        shutil.move(movepath, path)


def main():

    parser = argparse.ArgumentParser(description="Cleaner of the workflow")
    parser.add_argument('-m','--mapping', action="store_true", help="Just clean mapping columns")
    parser.add_argument('-a', '--all', action="store_true", help = "Clean all workflow outputs")
    args = parser.parse_args()

    if args.mapping:
        clean_mapping()
    if args.all:
        clean_total()


if __name__ == "__main__":
    main()