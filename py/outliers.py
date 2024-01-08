import logging 
import sys
from datetime import timedelta
import redis

from concurrent import futures

import grpc
import numpy as np

from outliers_pb2 import OutliersResponse
from outliers_pb2_grpc import OutliersServicer, add_OutliersServicer_to_server

def find_outliers(data: np.ndarray):
    out = np.where(np.abs(data - data.mean()) > 2 * data.std())
    return out[0]

class OutliersService(OutliersServicer):

    def __init__(self):
        self.r = redis.Redis(host='my-release-redis-master.default.svc.cluster.local', port=6379, password='7dnFhzZajA')
        self.pubsub = self.r.pubsub()
        self.pubsub.psubscribe(**{'__keyevent@0__:expired':self.key_event_handler})

    def Detect(self, request, context):
        logging.info('detect request size: %d', len(request.metrics))
        data = np.fromiter((m.value for m in request.metrics), dtype='float64')
        indicies = find_outliers(data)
        logging.info('found %d outliers', len(indicies))
        
        self.r.setex(str(np.random.default_rng().integers(1)), timedelta(seconds=5), "I have expired")
        resp = OutliersResponse(indicies=indicies)
        return resp
    
    def key_event_handler(self, msg):
        try:
            
            key = msg["data"].decode("utf-8")
            if "valKey" in key:
                key = key.replace("valKey:", "")
                value = self.r.get(key)
                logging.info('value of expired key: %d', value)
        except Exception as exp:
            pass
    
if __name__ == "__main__":
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(levelname)s - %(message)s',
        stream=sys.stdout,
    )
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    add_OutliersServicer_to_server(OutliersService(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    logging.info('server ready on port %r', 50051)
    server.wait_for_termination()