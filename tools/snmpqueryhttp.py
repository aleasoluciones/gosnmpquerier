# -*- coding: utf-8 -*-

import requests

data = {
    'cmd': 'get',
    'destination': 'localhost',
    'community': 'kaleidos',
    'oid': '1.3.6.1.2.1.2.2.1.10.1'
}

r = requests.get('http://localhost:8080/', params=data)

print('RESULTADO: ', r.text)
