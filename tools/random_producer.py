#!/usr/bin/env python
import puka

import argparse
import random




def generate_json_command(community):

	commands = (
		('walk', '1.3.6.1.2.1.2.2.1.2'),
		('walk', '1.3.6.1.2.1.2.2.1.10'),
		('get', '1.3.6.1.2.1.2.2.1.2.1'),
		('get', '1.3.6.1.2.1.1.5'),
	)
	destinations = (
		'ada-xem1', 'ona-xem1', 'alo-xem1', 'otm-xem1',
		'pob-xem1', 'c2k-xem1', 'vtr-xem1', 'tom-xem1',
		'onr-xem1', 'vco-xem1', 'inm-xem1', 'gtr-xem1',
		'ram-xem1', 'vir-xem1', 'tge-xem1', 'ola-xem1',
		'pip-xem1', 'vmc-xem1', 'pra-xem1', 'arm-xem1',
	)

	cmd, oid = random.choice(commands)
	destination = random.choice(destinations)
	json_format_string = '{"Command":"%s", "Destination":"%s.alea-soluciones.com", "Community":"%s", "Oid":"%s"}'
	return json_format_string % (cmd, destination, community, oid)
	
def main():
	parser = argparse.ArgumentParser()
	parser.add_argument('--amqp_uri', action="store", dest="amqp_uri")
	parser.add_argument('--community', action="store", dest="community")
	args = parser.parse_args()
	
	client = puka.Client(args.amqp_uri)
	promise = client.connect()
	client.wait(promise)
	for num in xrange(0,100):
		promise = client.basic_publish(
			exchange='EFA_SRC', 
			routing_key='', 
			body=generate_json_command(args.community))
		client.wait(promise)

	promise = client.close()
	client.wait(promise)

if __name__ == '__main__':
	main()
