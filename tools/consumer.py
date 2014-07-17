#!/usr/bin/env python
import puka
import json
import argparse

def main():
	parser = argparse.ArgumentParser()
	parser.add_argument('--amqp_uri', action="store", dest="amqp_uri")
	args = parser.parse_args()

	client = puka.Client(args.amqp_uri)
	promise = client.connect()
	client.wait(promise)

	
	consume_promise = client.basic_consume(queue='EFA_DST')
	try:
		while True:
		    result = client.wait(consume_promise)
		    print " [x] Received message %r" % (result,)
		    client.basic_ack(result)
	except Exception as exc:
		print "Error", exc, exc.__class__.__name__

	promise = client.basic_cancel(consume_promise)
	client.wait(promise)

	promise = client.close()
	client.wait(promise)

if __name__ == '__main__':
	main()