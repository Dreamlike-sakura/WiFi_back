import numpy as np
from Bfee import Bfee
from get_scale_csi import get_scale_csi
import scipy.io as io
import matplotlib.pyplot as plt
from scipy import signal
import json
import math as math
import sys as sys


#def read_bfee_file(datpath):
if __name__ == '__main__':
    bfee = Bfee.from_file(sys.argv[1], model_name_encode="gb2312")
    path=sys.argv[1]
    print(path)
    length=len(bfee.all_csi)

    data=np.array([ [complex(0, 0) for i in range(length)] for i in range(30) ])

    amp= np.array([[0.0 for i in range(length)] for i in range(30)])
    phase = np.array([[0.0 for i in range(length)] for i in range(30)])
    amp_dealt = np.array([[0 for i in range(length)] for i in range(30)])

    abnoraml1_w = np.array([0.0 for i in range(length)])
    abnoraml1_x = np.array([0.0 for i in range(length)])
    C = np.array([0 for i in range(length)])
    temparray= np.array([0.0 for i in range(30)])


    for j in range(0,len(bfee.all_csi)):
        csi = get_scale_csi(bfee.dicts[j]);  #numpy.complex128

        for i in range(30):
            data[i][j]=csi[i,0,0];


    for i in range(len(bfee.all_csi)):
        for j in range(30):
             amp[j][i]=abs(data[j][i])
             phase[j][i] = np.angle(data[j][i])
             b, a = signal.butter(2, 0.64, 'low')
             amp_dealt = signal.filtfilt(b, a, amp)

    data1=np.array(data)


    amp=amp.tolist()
    phase=phase.tolist()
    fp = open(path + '_amp.json', 'w')
    json_array1 = json.dump(amp, fp)
    fp = open(path + '_phase.json', 'w')
    json_array1 = json.dump(phase, fp)


