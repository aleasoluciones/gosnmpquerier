#!/usr/bin/env python
import puka

import argparse
import random
import data



def generate_json_command(community):
	cmd, oid = random.choice(data.commands)
	destination = random.choice(data.destinations)
	json_format_string = '{"Command":"%s", "Destination":"%s", "Community":"%s", "Oids":["%s"]}'
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
