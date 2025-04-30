import shutil

def main():

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
if __name__ == "__main__":
    main()