import time
import os, glob
import pandas as pd

path = "./Distribution"
path2 = "./Trend"
outpath = r'./Final Workbook/'

all_files = glob.glob(os.path.join(path, "*.csv"))
all_files2 = glob.glob(os.path.join(path2, "*.csv"))

df_from_each_file = (pd.read_csv(f, sep=',') for f in all_files)
df_from_each_file2 = (pd.read_csv(f, sep=',') for f in all_files2)

df_merged   = pd.concat(df_from_each_file, ignore_index=True)
df_merged2 = pd.concat(df_from_each_file2, ignore_index=True)

df_merged.to_csv(outpath + "Leadtime Distribution"+".csv",index=False)
df_merged2.to_csv(outpath + "Leadtime Trend"+".csv",index=False)

#df_merged.to_csv(outpath + "AAAA_"+(time.strftime("%Y-%m-%d")+".csv"))