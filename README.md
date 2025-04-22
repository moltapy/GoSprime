# GoSprime

<a href= "./README.en.md">English</a>

> Go语言封装的VCF处理工具(依赖Bcftools), 通过Python串联整个SPrime分析管线

## 背景

笔者是一位前生命科学专业的学生，因为课题分析需要接触到了SPrime这个分析工具，从整体上看整个分析的核心步骤是Browning团队写的`Sprime.jar`, 以及YingZhou001前辈用Shell、C和R脚本简单的串联了一整套流程

然而在我仍然在实验室的时期，初次接触整个软件的时候，有如下一些使用上的问题:

- YingZhou001前辈写的流程更多是适用于PBS集群的，对于一些没有部署PBS集群管理系统的集群适配性不好，需要学生自己改动代码
- Shell脚本串联的流程中，没有详细的日志和错误处理，让分析新手很难迅速找到错误的问题所在，比如数据格式导致的问题，笔者之前曾经在自己上手和同门的使用过程中遇到多次格式的问题，包括:
    - VCF文件染色体列是`chr{chrom}`导致在summary.r中作为number读入时报错
    - VCF文件的INFO列过长导致在Mapping步骤时, `DP[128]`数组越界而导致Mapping结果异常等
- `Sprime.jar`中含有大量无用的`ForkJoinPool`隐式线程,直接并行会导致整个集群异常的CPU资源消耗, 在执行时通过指定JVM参数`-Djava.util.concurrent.ForkJoinPool.common.parallelism`能够在不影响性能的情况下降低CPU的异常消耗

笔者也知道生物软件中很多轮子是特异的，这套流程在笔者写作文档时已经自测过了，有了相对完善的日志和日志分级，然而仍然无法保证能够兼容所有的数据和集群的情况。如果您在做分析时有幸或者不幸遇到了SPrime这款软件，又在其中碰见了一些自己不太好解决的问题，欢迎一起交流讨论，笔者虽然水平有限，仍然希望能尽己所能，提供一些力所能及的助力

在整个`GoSprime`项目的编写中，选择Go这门语言的原因有二:

- 其一是Go的语法比较接近数据分析常用的语言Python，相对简单易读
- 其二是Go语言天生高并发的优势，此外包体小依赖管理和编译方便

编写中R脚本部分直接使用了YingZhou001前辈的成果并进行了些微的改动，在此表示感谢，相关论文的引用和项目源码的传送门会放在文档的末尾

## 使用
整个SPrime分析流程由6个步骤组成:

- 从样本列表中按提取子群体，并将子群体和指定的外类群合并
- 根据子群体的样本，从VCF文件中提取对应样本的变异数据
- 按照染色体顺序，合并每个子群体的变异数据
- SPrime计算每个群体的分数，找到渗入的片段中的位点
- 将找到的渗入的片段位点和古人基因组进行匹配，找到匹配的位点
- 总结匹配古人基因组的渗入片段并根据总结数据作图

针对上述6个步骤，分别编写了如下的应用来处理:
- `ExtractSamples`:     负责从样本数据中提取子群体
- `ExtractAutos`:       根据提供的样本，从VCF文件中提取变异数据
- `ConcatAutos`:        合并每一个群体分染色体的变异数据文件
- `GenerateScores`:     使用SPrime推断出强连锁的片段中的位点
- `MappingArchs`:       将SPrime处理的结果位点和古人类的基因组匹配
- `Rscripts/`:           文件夹，内含两个脚本，分别负责总结匹配的结果和作图

每一个应用都可以单独作为一个工具使用，也是整体流程的一部分，下面会从整体到局部的角度介绍整个流程和所有工具

### 配置文件

整个`GoSprime`分析流程的默认配置文件在项目目录中的`config.yaml`中，也可以在运行流程的主要文件`sprime.py`时通过指定`-c`参数来指定其他的配置文件，配置文件中每一项的具体说明见配置文件的注释，这里不再赘述

### 编译运行

整个项目会依赖如下的内容: Go编译器、Bcftools、sprime.jar，下面依次讲解:

#### Go语言环境配置:

**Windows**
- 进入Go的官网 https://golang.google.cn/dl/ 下载安装包
- 在环境变量中配置`$GOROOT`和`$GOPATH`变量，分别指向Go的安装位置和工作位置
- 根据Go的版本选择配置`GOPROXY`和是否启用`GO MODULE`，参考Linux

**Linux(Ubuntu)**
```shell
# 进入Go官网，curl或者wget下载Go语言压缩文件
wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz

# 解压下载的压缩文件
tar -zvxf go1.24.2.linux-amd64.tar.gz

# 打开个人的bash配置
vim ~/.bashrc

# 在文件末尾增加如下配置(在用户的home目录下载解压的情况)
export GOROOT="$HOME/go" 
export PATH="$PATH:$GOROOT/bin"
export GOPATH=$HOME/go/lib:$HOME/go/work

# 验证安装
go version  # 这一步显示版本号说明安装成功

# 如果有网络问题，配置 GOPROXY 环境变量，以下三选一
## 1. 七牛云
go env -w GOPROXY=https://goproxy.cn,direct
## 2. 阿里云
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
## 3. 官方
go env -w GOPROXY=https://goproxy.io,direct
# 确认设置
go env | grep GOPROXY

# 如果Go版本低于1.1.2， 需要手动启用Go Modules
go env -w GO111MODULE=on
```

#### Bcftools

`Bcftools`工具的安装常见的有两种方式：
- 如果有sudo权限，并且确定不会污染集群内环境：`sudo apt install bcftools`，配置文件中bcftools部分可以留空
- 如果使用conda做环境变量管理，那么如果安装的是非base环境，由于切换时会导致路径变化，因此只能配置bcftools的路径
    ```shell
    # 开启一个新的环境名字叫sprime
    conda create -n sprime
    # 通过bioconda channel安装bcftools
    conda install -c bioconda bcftools
    ```
    此时bcftools工具的路径：`~/miniconda3/envs/sprime/bin/bcftools`
#### Java环境配置

运行`sprime.jar`需要配置Java环境，这里配置使用的是Java 17

```shell
# 创建Java目录，下载Java二进制文件夹并解压
mkdir ~/java
cd ~/java
wget https://download.oracle.com/java/17/archive/jdk-17.0.10_linux-aarch64_bin.tar.gz
tar -zvxf jdk-17.0.10_linux-aarch64_bin.tar.gz
# 将Java配置进入环境变量并刷新bash使其生效
vim ~/.bashrc
export JAVA_HOME=~/java/jdk-17.0.10_linux-aarch64_bin   
export PATH=$PATH:$JAVA_HOME/bin
source ~/.bashrc

# 验证Java安装的结果
java -version
```


在完成上述配置之后(如果不需要SPrime处理的步骤可以省略Java环境的配置),有以下两种情况：
- 如果仅需使用其中部分程序，进入程序文件夹，使用`go build`编译，程序会自动处理相关依赖的下载
- 如果需要运行流程，在配置好`config.yaml`后，运行`python sprime.py`, Python会处理程序构建

### 详细说明

#### ExtractSamples

`ExtractSamples`应用的参数如下：
- `-h`, `--help`           : 打印整个程序的帮助信息并退出程序
- `-s`, `--sample`         : 样本信息文本文件，第一列是样本ID，第二列是样本所属的子群体，默认应该是空格或者`\t`作为分隔符，每一列的分隔符要统一
- `-g`, `--outgroup`       : 外类群的群体名称
- `-p`, `--samplepath`     : 样本子群体的输出路径，输出的子群体会以`<groupname>.txt`的方式存储在其中，如果指定外类群,结果会合并子群体和外类群的样本，如果不指定，可以分离提取所有子群体的样本
- `-o`, `--outputgrouppath`: 外类群样本的输出路径，外类群会以`outgroup.txt`的命名存储在其中
- `-l`, `--limits`         : 群体的范围，以`,`分隔的方式输入提取的群体，只提取范围内的群体，默认使用全部群体
- `-t`, `--title`          : 标记是否包含一行标头，默认值是`true`
- `-f`, `--format`         : 用于格式比较混乱的情况，用`|`代表每一行输入解析的字符串，能够成功解析，比如`ID<space>Subgroup\tGroup<space>\tgender`,输入`|<space>|\t|<space>\t|`

#### ExtractAutos

`ExtractAutos`应用的参数如下：
- `-h`, `--help`           : 打印整个程序的帮助信息并退出程序
- `-s`, `--samplepath`     : 群体样本文本文件的保存目录，每个群体的文本文件都只有一列样本ID，路径是其父目录
- `-m`, `--modernpath`     : 现代人类VCF文件的路径，路径需要具体到每个文件，用`{chr}`作为染色体号的占位符
- `-o`, `--modernout`      : 提取每个群体样本变异信息的输出目录，每个群体的染色体变异信息会保存在对应的`<groupname>`目录下
- `-n`, `--name`           : 提取后每个群体每条染色体对应VCF文件的名字，用`{chr}`最为染色体号的占位符，默认值是`chrom_{chr}.vcf.gz`
- `-a`, `--autos`          : 染色体号的范围，以`,`分隔的方式输入提取的染色体号，只提取范围内的染色体，默认使用全部染色体
- `-g`, `--groups`         : 群体的范围，以`,`分隔的方式输入提取的群体，只提取范围内的群体，默认使用全部群体
- `-b`, `--bcftools`       : `Bcftools`工具的路径，如果全局安装可以留空字符串，如果安装在其他路径，提供绝对路径
- `-t`, `--threads`        : 指定在并发运行时候的线程数，默认值根据集群或者机器核数的不同而有差异，默认为核数的1/2
- `-c`, `--cores`          : 指定在并发运行时使用的核数，默认值根据集群或者机器核数的不同而有差异，在16核以上是1/4的核；16核以下最多使用8核
- `--overwrite`            : 标记是否允许覆盖，如果不允许覆盖，检测到输出目录存在且非空，会直接使用该目录的结果进行下一步计算

#### ConcatAutos

`ConcatAutos`应用的参数如下：
- `-h`, `--help`           : 打印整个程序的帮助信息并退出程序
- `-d`, `--dir`            : 样本按照子群体提取后生成单个群体VCF文件的目录，是包含`<groupname>`目录的父目录
- `-g`, `--group`          : 群体的范围，以`,`分隔的方式输入用于合并处理的群体，只合并范围内群体的染色体变异信息，默认使用全部群体
- `-o`, `--outdir`         : 每个群体合并后VCF文件的输出目录，每个群体合并后的文件将保存为`<groupname>_all.vcf.gz`
- `-b`, `--bcftools`       : `Bcftools`工具的路径，如果全局安装可以留空字符串，如果安装在其他路径，提供绝对路径
- `--overwrite`            : 标记是否允许覆盖，如果不允许覆盖，检测到输出目录存在且非空，会直接使用该目录的结果进行下一步计算

#### GenerateScores

`GenerateScores`应用的参数如下：
- `-h`,`--help`            : 打印整个程序的帮助信息并退出程序
- `-g`,`--genopath`        : 合并后每个群体VCF文件的路径，群体需要使用`{group}`作为占位符，如上述默认为`path/to/{group}_all.vcf.gz`
- `-m`, `--maproute`       : PLINK生成的参考基因组的map文件路径
- `-j`, `--jarpath`        : `sprime.jar`文件的路径
- `-u`, `--outgroup`       : 记载外类群样本信息的文本文件路径，一列包含ID的文本文件，`ExtractSamples`会自动生成
- `-n`, `--outname`        : 外类群的群体名称
- `-o`, `--output`         : SPrime识别到每个群体连锁片段结果文件的目录，是包含`<groupname>`目录的父目录，每条染色体的结果以`.score`结尾
- `-r`, `--range`          : 群体的范围，以`,`分隔的方式输入用于合并处理的群体，只合并范围内群体的染色体变异信息，默认使用全部群体
- `-l`, `--chromlist`      : 染色体号的范围，以`,`分隔的方式输入提取的染色体号，只计算范围内的染色体，默认使用全部染色体
- `-t`, `--threads`        : 指定在并发运行时候的线程数，默认值根据集群或者机器核数的不同而有差异，默认为核数的1/2
- `-c`, `--cores`          : 指定在并发运行时使用的核数，默认值根据集群或者机器核数的不同而有差异，在16核以上是1/4的核；16核以下最多使用8核
- `--overwrite`            : 标记是否允许覆盖，如果不允许覆盖，检测到输出目录存在且非空，会直接使用该目录的结果进行下一步计算

#### MappingArchs

`MappingArchs`应用的参数如下：
- `-h`, `--help`           : 打印整个程序的帮助信息并退出程序
- `-m`, `--maskpath`       : 移除掩码文件的路径，gzip压缩的BED文件，染色体序号使用`{chr}`占位
- `-v`, `--vcfarch`        : 古人VCF文件的路径，gzip压缩的VCF文件，染色体需要使用`{chr}`占位
- `-p`, `--scorepath`      : 每个群体SPrime结果文件的父目录，是包含`<groupname>`目录的父目录，每条染色体的结果以`.score`结尾
- `-n`, `--arrayname`      : 增加和群体匹配结果列的列名
- `-s`, `--separator`      : 输出增加列时所使用的分隔符
- `--reverse`              : 提供了`maskpath`参数后可用，表示提供的BED文件中的位点需要包含而非排除
- `--depth`                : 在增加的匹配结果列后额外增加一列表示测序深度的列

    > 注: 如果在Windows上使用，注意可能`os.Rename()`原子替换函数可能因为系统的原因不成功

## 证书

根据项目所使用的部分开源库的要求，使用Apache证书开源

## 引用

若您需要快速查找相关引用文献，列举如下，相关项目原始传送门也一并附上:

Sprime项目地址: https://github.com/browning-lab/sprime
> SPrime: S R Browning, B L Browning, Y Zhou, S Tucci, J M Akey (2018). Analysis of human sequence data reveals two pulses of archaic Denisovan admixture. Cell 173(1):53-61. doi: 10.1016/j.cell.2018.02.031

Sprime流程文档项目地址: https://github.com/YingZhou001/sprimepipeline
> Pipeline: Y Zhou, S R Browning (2021). Protocol for detecting archaic introgression variants with SPrime. Submitted.
