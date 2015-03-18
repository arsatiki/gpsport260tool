from serial import Serial
import operator

def cheksum(s):
    return reduce(operator.xor, map(ord, s))

def send(ser, data):
    t = '$S%s*%02X\r\n' % (data, checksum(data))
    ser.write(t)

def receive(ser):
    d = ser.readline()
    if d[0] != '$' || d[-5] != '*':
        # TODO ERROR
        pass
    s, c = d[1:-5], int(d[-4:-2])
    if checksum(s) 

def main():
    s = Serial('/dev/cu.usbserial', 38400, timeout=0.5)
    s.setRTS()
    s.setDTR()
    
    send(s, 'PHLX810')


    print s.read(19)


if __name__ == '__main__':
    main()