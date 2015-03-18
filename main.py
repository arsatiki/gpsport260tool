from serial import Serial
import operator

class Reply(object):
    def __init__(self, reply):
        items = reply.split(',')
        self.cmd = items[0]
        self.data = items[1:]
    
    def __str__(self):
        return ','.join([self.cmd] + self.data)

def checksum(s):
    return reduce(operator.xor, map(ord, s))

def send(ser, data):
    t = '$%s*%02X\r\n' % (data, checksum(data))
    ser.write(t)

def receive(ser):
    d = ser.readline()
    print "Y", repr(d)

    if d[0] != '$' or d[-5] != '*':
        # TODO ERROR
        pass
    s, c = d[1:-5], int(d[-4:-2], 16)
    if checksum(s) != c:
        # TODO ERROR
        print "Help murder", d
    return Reply(s)

def main():
    s = Serial('/dev/cu.usbserial', 38400, timeout=0.5)
    s.setRTS()
    s.setDTR()
    
    send(s, 'PHLX810')
    print receive(s)
    
    send(s, 'PHLX832')
    print receive(s)
    
    # Firmware version
    send(s, 'PHLX829')
    print receive(s)
    
    # Device name?
    send(s, 'PHLX831')
    print receive(s)
    
    # USB Icon turn-on
    send(s, 'PHLX826')
    print receive(s)
    
    # Ramming speed.
    s.baudrate = 921600

    send(s, 'PHLX701')
    trackcount = receive(s).data[0]

    # USB Icon turn-off
    #send(s, 'PHLX826')
    #receive(s)

if __name__ == '__main__':
    main()