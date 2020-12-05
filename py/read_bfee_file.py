import numpy as np
from Bfee import Bfee
from get_scale_csi import get_scale_csi
import scipy.io as io
import matplotlib.pyplot as plt
from scipy import signal
import json
import math as math
import sys as sys
pi=3.14

#def read_bfee_file(datpath):
if __name__ == '__main__':

    filename=sys.argv[1]
    user=sys.argv[2]
    c = filename[0]
    type = 'run'
    if (c == 'r'):
        type = 'run'
    elif (c == 'w'):
        type = 'walk'
    elif (c == 's'):
        type = 'shake'

    upload='../data/'+user+'/upload/'+type+'/'+filename+'.dat'
    bfee = Bfee.from_file(upload, model_name_encode="gb2312")
    # fore_path='D:\20study\2020project\back\data\dealt'


    originamppath='../data/'+user+'/origin/'+type+'/amp/'
    originphasepath='../data/'+user+'/origin/'+type+'/phase/'
    dealtamppath='../data/'+user+'/dealt/'+type+'/amp/'
    dealtphasepath='../data/'+user+'/dealt/'+type+'/phase/'
    abnormalpath='../data/'+user+'/dealt/'+type+'/abnormal/'

    length=len(bfee.all_csi)

    data=np.array([ [complex(0, 0) for i in range(length)] for i in range(30) ])

    amp= np.array([[0.0 for i in range(length)] for i in range(30)])
    phase = np.array([[0.0 for i in range(length)] for i in range(30)])
    amp_dealt = np.array([[0 for i in range(length)] for i in range(30)])
    phase_dealt = np.array([[0 for i in range(length)] for i in range(30)])
    m = [-28, -26, -24, -22, -20, -18, -16, -14, -12, -10, -8, -6, -4, -2, -1, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21,
         23, 25, 27, 28]

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

    temparray = np.array([[0 for i in range(length)] for i in range(30)])
    n = pi
    for i in range(length):
        diff = 0
        temparray[0][i] = phase[0][i]
        for j in range(1, 30):
            if ((phase[j][i] - phase[j - 1][i]) > pi):
                diff = diff + 1
            temparray[j][i] = phase[j][i] - diff * pi * 2

    sum = np.array([0 for i in range(length)])
    for i in range(length):
        for j in range(30):
            sum[i] += phase[j][i]
    for i in range(length):
        k = (temparray[29][i] - temparray[0][i]) / (m[29] - m[0])
        for j in range(30):
            phase_dealt[j][i] = temparray[j][i] - k * m[j] - sum[i] / 30

    data1 = np.array(data)
    limit2=50

    mat = np.array(data1);
    amp = np.array(amp);
    phase = np.array(phase);

    for i in range(length):
        temparray = (amp[:, i])
        abnoraml1_w[i] = np.var(temparray)

    e = np.mean(abnoraml1_w)
    w = np.var(abnoraml1_w)
    sw = math.sqrt(w)
    for i in range(length):
        abnoraml1_x[i] = (abnoraml1_w[i] - e) / sw
    limit1 = 1.5 * (np.std(abnoraml1_x))

    for i in range(length):
        if (abnoraml1_w[i] >= limit1):
            C[i] = 1;
        else:
            C[i] = 0

    me = np.array([0 for i in range(length)])
    w = np.array([0 for i in range(length)])
    sum1 = np.array([0 for i in range(length)])
    me = np.mean(data, axis=0)
    w = np.var(data, axis=0)
    sum1 = np.sum(amp, axis=0)

    for i in range(length):
        if (sum1[i] >= limit2):
            sum1[i] = 1
        else:
            sum1[i] = 0

    start_point = np.array([0 for i in range(length)])
    for i in range(length - 1):
        if (sum1[i] == 1 and sum1[i - 1] == 0):
            start_point[i] = i
        if (sum1[i] == 0 and sum1[i + 1] == 0):
            start_point[i] = 1

    amp=amp.tolist()
    phase=phase.tolist()

    fp = open(originamppath+filename + '_amp.json', 'w')
    json_array1 = json.dump(amp, fp)
    fp = open(originphasepath+filename+ '_phase.json', 'w')
    json_array1 = json.dump(phase, fp)

    amp_dealt=amp_dealt.tolist()
    phase_dealt=phase_dealt.tolist()

    fp = open(dealtamppath+'dealt_'+filename+'_amp.json', 'w')
    json_array1 = json.dump(amp_dealt, fp)
    fp = open(dealtphasepath+'dealt_'+filename+'_phase.json', 'w')
    json_array1 = json.dump(phase_dealt, fp)

    start_point = start_point.tolist()
    start_point=[start_point for i in range(1)]

    
    fp = open(abnormalpath +'abnormal_' +filename+'.json', 'w')
    json_array1 = json.dump(start_point, fp)


