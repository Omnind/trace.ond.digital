import os, glob
import time
import pandas as pd

file_path='./Distribution'
outpath = './Trend'

file_path=file_path.strip('/')
file_name_list= os.listdir(file_path)
for index,file_name in enumerate(file_name_list):
    if '.csv' not in file_name:
        continue
    df=pd.read_csv(file_path+'/'+file_name, low_memory=False)

    df = df[df['Start Station'] == "CNC Input"]
    df = df[df['End Station'] == "FQC"]
    df=df.drop(['Start Station','End Station', 'Parts','Average'],axis=1)

    df.to_csv(outpath+'/'+file_name,index=False)
