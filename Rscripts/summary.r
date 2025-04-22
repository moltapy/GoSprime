args = commandArgs(trailingOnly=TRUE)

calmatchrate<-function(X, minn)
{
    total<-length(grep("match", X))
    mis<-length(grep("mismatch", X))
    if(total>=minn)return(round(1-mis/total, digits=4))
    else return(NA)
}

minn=10

files=list.files(path = args[1], pattern = "*.score$")
files<-paste(args[1],"/",files,sep="")

out<-c()
for(score in files){
    da<-read.csv(file=score, header=T,sep="\t")
    chr=as.character(da[1,1])
    seglist<-unique(da$SEGMENT)
    for(i in 1:length(seglist)){
        outlist<-da$SEGMENT==seglist[i]
        x<-c(chr, seglist[i], range(da$POS[outlist]))
        y<-c()
        for(j in 9:dim(da)[2]){
            y<-c(y, calmatchrate(da[outlist, j], minn))
        }
        out<-rbind(out, (c(x,y)))
    }
}

colnames(out)<-c("chr", "seg", "from", "to", colnames(da)[9:dim(da)[2]])
out<-out[order(out[,1]),]

write.table(out, args[2], quote=F, row.names=F)