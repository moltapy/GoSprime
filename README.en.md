# GoSprime

<a href="./README.md">简体中文</a>

> End-to-end SPrime workflow automation using Python with companion VCF processing tools 

## Background

As a former life sciences student, I encountered the SPrime analysis tool during my research project. The core workflow primarily consists of two components:  
1) The Sprime.jar package developed by the Browning team, which handles the central computational algorithms.  
2) A series of Shell, C, and R scripts created by researcher YingZhou001 to orchestrate the pipeline.  

During my initial implementation in the lab environment, several usability challenges emerged:  

1. Cluster Compatibility Issues  
    The original workflow was optimized for PBS cluster systems, requiring manual code modifications for non-PBS HPC environments.  


2. Pipeline Debugging Limitations  
    The shell script-based workflow lacks detailed logging and error handling mechanisms, making it difficult for novice users to quickly pinpoint issues. Specific data format challenges encountered include:  

    - VCF files with chromosome columns formatted as `chr{chrom}` caused errors in the `summary.r` script when parsed as numeric values  

    - INFO field truncation issues in VCF files leading to array boundary violations (e.g., `DP[128]` overflow) during mapping steps  
    ......

3. Resource Management Optimization  
    The Sprime.jar implementation contains redundant ForkJoinPool implicit threading, which caused abnormal CPU consumption in cluster environments. We mitigated this through JVM parameter tuning:  
  `-Djava.util.concurrent.ForkJoinPool.common.parallelism=[custom_thread_num]`
  This adjustment maintained computational efficiency while reducing resource contention.  

While I've implemented enhanced logging (with severity levels) and validated the workflow through self-testing during documentation development, complete compatibility across all data types and cluster configurations cannot be guaranteed.  

If you encounter challenges when implementing the SPrime toolchain – whether through fortunate adoption or unavoidable necessity – feel free to reach out for discussion. Though my expertise is limited, I will endeavor to provide whatever assistance I can within my capabilities.  

The GoSprime project adopted the Go programming language for two primary reasons:  
1. Syntax Familiarity: Go's syntax shares similarities with Python (widely used in data analysis), ensuring code readability and a gentle learning curve.  
2. Technical Advantages:  
   - Native support for high-concurrency paradigms  

   - Compact binaries with streamlined dependency management  

   - Straightforward cross-platform compilation  

The R scripting components directly utilize YingZhou001's original implementation, with minor modifications to enhance compatibility. I gratefully acknowledge these foundational contributions. Full citations to relevant publications and a portal to the source code repository will be provided in the References section of this documentation.  


## Usage

The SPrime analytical workflow consists of six sequential stages:  

1. Subpopulation Extraction:  
   Extract target subpopulations from the sample list and merge them with designated outgroups.  

2. Variant Data Extraction:  
   Retrieve variant data for specified subpopulations from VCF files.  

3. Chromosomal Data Consolidation:  
    Merge subpopulation-specific variant data by chromosome.  

4. Introgression Site Detection:  
   Calculate population-specific scores using SPrime to identify loci within strongly linked introgressed segments.  

5. Ancient Genome Matching:  
   Cross-reference detected loci with ancient human genomic data.  

6. Result Visualization:  
   Summarize matched introgressed segments and generate graphical representations.  

Each stage corresponds to a dedicated executable tool:  

| Component           | Functionality |  
|---------------------|---------------|  
| `ExtractSamples`    | Handles subpopulation extraction from sample data |  
| `ExtractAutos`      | Extracts chromosome-specific variants from VCF files |  
| `ConcatAutos`       | Merges variant data across chromosomes per subpopulation |  
| `GenerateScores`    | Identifies candidate loci using SPrime scoring algorithms |  
| `MappingArchs`      | Performs ancient genome alignment for detected loci |  
| `Rscipts/`          | Contains two R scripts:  
|                     | - `summary.r`: Statistical aggregation of matched |segments  
|                     | - `plotting.r`: Visualization of final results  |

Each tool functions both as:  
- A standalone utility for specific analytical tasks  
- An integrated component within the complete pipeline  

The following sections will provide detailed documentation, progressing from holistic workflow explanations to individual tool specifications.

### Configuration File  

The default configuration file for the `GoSprime` pipeline resides at `config.yaml` in the project directory. Users may specify alternative configurations using the `-c` parameter when executing the main script `sprime.py`. Detailed parameter explanations are provided in the configuration file comments.

### Compilation & Execution  

**Dependencies**  
- Go compiler (≥1.12 recommended)  
- Bcftools (≥1.10)  
- sprime.jar

#### Go Environment Setup  

**Windows**  
1. Download installer from [Go Official Site](https://golang.google.cn/dl/)  
2. Configure environment variables:  
   ```text
   GOROOT=[Go installation path]  
   GOPATH=[Workspace path]  
   ```  

**Linux (Ubuntu)**

```shell
wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz
tar -zvxf go1.24.2.linux-amd64.tar.gz

# Add to ~/.bashrc
export GOROOT="$HOME/go"
export PATH="$PATH:$GOROOT/bin"
export GOPATH="$HOME/go/lib:$HOME/go/work"

# Verify installation
go version
```

---

#### Bcftools Installation  

System-wide (requires sudo):  
```shell
sudo apt install bcftools
```

Conda Environment:  
```shell
conda create -n sprime
conda install -c bioconda bcftools
# Path: ~/miniconda3/envs/sprime/bin/bcftools
```

---

#### Java Setup(Example: Java 17)  

```shell
mkdir ~/java && cd ~/java
wget https://download.oracle.com/java/17/archive/jdk-17.0.10_linux-aarch64_bin.tar.gz
tar -zvxf jdk-17.0.10_linux-aarch64_bin.tar.gz

# Add to ~/.bashrc
export JAVA_HOME=~/java/jdk-17.0.10_linux-aarch64_bin
export PATH="$PATH:$JAVA_HOME/bin"

# Verify
java -version
```
---

#### Execution Modes  
1. Component-wise Compilation  
   ```shell
   cd [component_directory] && go build
   ```  
2. Full Pipeline  
   ```shell 
   python sprime.py -c [config_path]
   ```

---

### Tool Specifications  

**ExtractSamples**  

| Option & Shorthand       | Description |  
|--------------------------|-------------|  
| `-h`, `--help`           | Display help message and exit |  
| `-s`, `--sample`         | Sample metadata text file containing two columns:<br>1) Sample IDs<br>2) Corresponding subpopulation<br>Format: Space or `\t` delimited (consistent per line) |  
| `-g`, `--outgroup`       | Name of outgroup population |  
| `-p`, `--samplepath`     | Output directory for subpopulation lists:<br>- Creates `<groupname>.txt` per subpopulation<br>- Merges subpopulations with outgroup if specified |  
| `-o`, `--outputgrouppath`| Output path for outgroup samples:<br>- Generates `outgroup.txt` containing merged outgroup IDs |  
| `-l`, `--limits`         | Comma-separated list of target subpopulations to extract<br>Default: Process all groups |  
| `-t`, `--title`          | Indicates header presence in input file<br>Values: `true`/`false`<br>Default: `true` |  
| `-f`, `--format`         | Custom parser for inconsistent formatting:<br>- Use `\|` to define column separators<br>Example: For input `ID<space>Subgroup\tGroup<space>\tgender`<br>Specify: `"\|<space>\|\t\|\<space>\t\|"` |  

--- 

**ExtractAutos**  

| Option & Shorthand       | Description |  
|--------------------------|-------------|  
| `-h`, `--help`           | Display help message and exit |  
| `-s`, `--samplepath`     | Parent directory containing subpopulation sample lists<br>- Each group's file: single-column sample IDs (e.g., `group1.txt`) |  
| `-m`, `--modernpath`     | Template path to modern human VCF files<br>- Use `{chr}` as chromosome placeholder (e.g., `data/chr{chr}.vcf.gz`) |  
| `-o`, `--modernout`      | Output directory for extracted variants<br>- Creates `<groupname>` subdirectories per population |  
| `-n`, `--name`           | Output VCF filename template<br>- Default: `chrom_{chr}.vcf.gz`<br>- `{chr}` represents chromosome number |  
| `-a`, `--autos`          | Comma-separated list of chromosomes to process<br>Default: All autosomes |  
| `-g`, `--groups`         | Comma-separated list of target subpopulations<br>Default: All groups |  
| `-b`, `--bcftools`       | Path to bcftools executable<br>- Empty string if globally installed<br>- Absolute path required for custom locations |  
| `-t`, `--threads`        | Number of concurrent threads<br>Default: 50% of total CPU cores |  
| `-c`, `--cores`          | Maximum CPU cores allocated:<br>- >16 cores: 25% of total<br>- ≤16 cores: max 8 cores |  
| `--overwrite`            | Overwrite existing output directories<br>- If disabled: skips processing when non-empty output exists |  

---

**ConcatAutos**  
| Option & Shorthand       | Description |  
|--------------------------|-------------|  
| `-h`, `--help`           | Display help message and exit |  
| `-d`, `--dir`            | Parent directory containing subpopulation-specific VCFs<br>- Expected structure: `[dir]/<groupname>/chrom_{chr}.vcf.gz` |  
| `-g`, `--group`          | Comma-separated list of subpopulations to process<br>Default: All groups |  
| `-o`, `--outdir`         | Output directory for merged VCFs<br>- Generates `<groupname>_all.vcf.gz` per subpopulation |  
| `-b`, `--bcftools`       | Path to bcftools executable<br>- Empty string if in system PATH<br>- Absolute path required for custom installations |  
| `--overwrite`            | Overwrite pre-existing output files<br>- If disabled: Skips processing when output file exists |  

---

**GenerateScores**  

| Option & Shorthand       | Description |  
|--------------------------|-------------|  
| `-h`, `--help`            | Display help message and exit |  
| `-g`, `--genopath`        | Path template to merged VCF files per subpopulation<br>- Use `{group}` placeholder (e.g., `data/{group}_all.vcf.gz`) |  
| `-m`, `--maproute`        | Path to PLINK-generated reference genome map file |  
| `-j`, `--jarpath`         | Absolute path to `sprime.jar` executable |  
| `-u`, `--outgroup`        | Single-column text file containing outgroup sample IDs<br>- Auto-generated by `ExtractSamples` |  
| `-n`, `--outname`         | Designated name for outgroup population |  
| `-o`, `--output`          | Output directory hierarchy:<br>- Creates `<groupname>` subdirectories<br>- Stores per-chromosome results as `.score` files |  
| `-r`, `--range`           | Comma-separated subpopulation filter<br>Default: Process all groups |  
| `-l`, `--chromlist`       | Comma-separated chromosome filter<br>Default: All chromosomes |  
| `-t`, `--threads`         | Thread concurrency level<br>Default: 50% of total CPU cores |  
| `-c`, `--cores`           | CPU core allocation policy:<br>- Systems with >16 cores: 25% allocated<br>- Systems with ≤16 cores: Max 8 cores |  
| `--overwrite`             | Override existing output files<br>- If disabled: Uses existing results when output directory contains valid files |  

---

**MappingArchs**

| Option & Shorthand       | Description |  
|--------------------------|-------------|  
| `-h`, `--help`           | Display help message and exit |  
| `-m`, `--maskpath`       | Path to gzipped BED mask file<br>- Chromosome placeholder: `{chr}` (e.g., `masks/chr{chr}.bed.gz`) |  
| `-v`, `--vcfarch`        | Path template to ancient genome VCF files<br>- Requires gzip compression<br>- Chromosome placeholder: `{chr}` |  
| `-p`, `--scorepath`      | Parent directory of SPrime results<br>- Expected structure: `[path]/<groupname>/chr*.score` |  
| `-n`, `--arrayname`      | Column header name for population matching results |  
| `-s`, `--separator`      | Delimiter for appended result columns<br>Default: Tab (`\t`) |  
| `--reverse`              | Invert mask file logic (include instead of exclude BED regions) |  
| `--depth`                | Add sequencing depth column to output |  

Platform Note:  
> When running on Windows systems, the `os.Rename()` atomic replacement function may fail due to OS-level file locking mechanisms. Consider implementing alternative file overwrite strategies.  

---

## Licensing  
This project is licensed under Apache License 2.0 to comply with dependency requirements.

## Citations  

SPrime Core Algorithm, GitHub: https://github.com/browning-lab/sprime    
> S R Browning, B L Browning, Y Zhou, S Tucci, J M Akey (2018). *Analysis of human sequence data reveals two pulses of archaic Denisovan admixture.* Cell 173(1):53-61. doi: 10.1016/j.cell.2018.02.031  


Pipeline Implementation, GitHub: https://github.com/YingZhou001/sprimepipeline  
> Zhou Y, Browning SR (2021). *Protocol for detecting archaic introgression variants with SPrime* .  
