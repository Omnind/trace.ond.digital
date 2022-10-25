import os
import pandas as pd

# Configure Folder Name
file_path='./demodata/large'
headers2=['root_serial', '2d-bc-le.local_time', '2d-bc-le.insight.test_attributes.uut_start', '2d-bc-le.insight.test_attributes.uut_stop']

file_path=file_path.strip('/')
file_name_list= os.listdir(file_path)
for index,file_name in enumerate(file_name_list):
    if '.csv' not in file_name:
        continue
    headers=[]
    qz=file_name.split('TI_')[-1].split('_')[0]
    df=pd.read_csv(file_path+'/'+file_name, low_memory=False)

    headers3=[x.split('.')[-1] for x in   df.columns.tolist()]
    for h in headers2:
        if h.split('.')[-1] in headers3:
            headers.append(df.columns.tolist()[headers3.index(h.split('.')[-1])])
        else:
            headers.append(h.replace('2d-bc-le',qz))

    for h in headers[-2:]:
        if h not in df.columns.tolist():
            df[h]=df['2d-bc-le.local_time'.replace('2d-bc-le',qz)]
    for h in headers:
        if h not in df.columns.tolist():
            df[h]=''


    df[headers].to_csv(file_path+'/'+file_name,index=False)
    print(index+1,len(file_name_list),file_name,'Done！')

