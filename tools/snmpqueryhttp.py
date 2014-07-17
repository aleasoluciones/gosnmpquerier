# -*- coding: utf-8 -*-

import requests
import data
import argparse
import random

def main():
	parser = argparse.ArgumentParser()
	parser.add_argument('--community', required=True, action="store", dest="community")
	args = parser.parse_args()

	for num in xrange(0,20):
		cmd, oid = random.choice(data.commands)
		destination = random.choice(data.destinations)
		query = {
			'cmd': cmd,
			'destination': destination,
			'community': args.community,
			'oid': oid
		}

		r = requests.get('http://localhost:8080/', params=query)
		print('RESULTADO: ', r.text)



if __name__ == '__main__':
	main()
