with open("test.score","wt") as ouf:
    for i in range(10):
        j =[str(x) for x in list(range(10,20))]
        ouf.write("\t".join(j)+"\n")
