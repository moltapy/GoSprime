#!bin/bash
conda activate bcftools
cd /public/group_data_2023/Heyuan/Sprime/output/subpops/GoSprime
go build -o splitGroup ./SplitGroup

cd ../..

subpops/GoSprime/splitGroup -m /public/group_data/he_yuan/data/sekei/data_1kg/modern_v5b/chrom_{chrom}.vcf.gz -o YRI -s /public/group_data/he_yuan/IBDmix_related/Samplelists/Original/sample_all.txtonly1000g -p 4



