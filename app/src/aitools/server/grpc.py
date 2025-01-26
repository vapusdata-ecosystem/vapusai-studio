from concurrent import futures
from loguru import logger
import sys
import os
# from grpc_reflection.v1alpha import reflection
from utils.importer import proto_importer
from interceptors.interceptor import Interceptor
from database.connector import DatabaseConnector
import argparse
proto_importer()


import grpc
from protos.vapus_aiutilities.v1alpha1 import vapus_aiutilities_pb2_grpc
import protos.vapus_aiutilities.v1alpha1.vapus_aiutilities_pb2 as pb2
from helpers.config import load_vapusaiserver_config,VapusAiConfig
# from helpers.secrets import init_vapus_backend_secrets, SecretStore
from helpers.logger import *
from helpers import settings
# from server.boot import ServerBoot
from controller.vapusmlutilities import AIUtilityService

class GrpcServer:
    """
    Represents a gRPC server for the VapusAi service.
    """

    serviceConfig: VapusAiConfig
    # secretsConfig: SecretStore
    
    @classmethod
    def configure_logger(cls,args:list):
        """
        Configures the logger based on the command line arguments.

        Args:
            args (list): The command line arguments.
        """
        logger.remove()
        if "--debug" in args:
            logger.add(sys.stderr, format="{time} {level} {message}", level="DEBUG")
            config = {
            "handlers": [
                {"sink": sys.stdout, "level": "DEBUG"},
                ]
            }
        else:
            logger.add(sys.stderr, format="{time} {level} {message}", level="INFO")
            config = {
            "handlers": [
                {"sink": sys.stdout, "level": "INFO"},
            ]
        }
        logger.configure(**config)
        logger.info("Logger configured")
        service_logger = logger

    def init_server(cls):
        """
        Initializes and starts the gRPC server.

        This method initializes a gRPC server, adds the VapusAiServiceServicer to the server,
        starts the server on the specified port, and waits for termination.

        Args:
            cls: The class object.

        Returns:
            None
        """
        port = cls.serviceConfig.networkConfig.mlutilitySvc.port
        service_logger.info("Starting server on port {port}", port=port)
        try:
            cls.server = grpc.server(futures.ThreadPoolExecutor(max_workers=10),interceptors=[Interceptor()])
            vapus_aiutilities_pb2_grpc.add_AIUtilityServicer_to_server(AIUtilityService(), cls.server)
            '''
            connecting to database
            '''
            

            # SERVICE_NAMES = (
            #     pb2.DESCRIPTOR.services_by_name['AIUtility'].full_name,
            #     reflection.SERVICE_NAME,
            #     )
            # reflection.enable_server_reflection(SERVICE_NAMES, cls.server)
            cls.server.add_insecure_port("[::]:" + str(port))
        except Exception as e:
            service_logger.error("Error initializing grpc server: {error}", error=str(e))
            raise e
        
       
    def startDBEngine(cls):
        db_connector = DatabaseConnector()
        dbSecret = cls.serviceConfig.mainConfig.vapusBESecretStorage.secret
        try:
            engine = db_connector.NewConnection(dbSecret)
        except Exception as e:
            raise(e)
        return engine
    def configure(cls):
        """
        Starts the gRPC server by loading the service configuration, initializing secrets, configuring the logger, and initializing the server.
        """ 
        
        parser = argparse.ArgumentParser(description="Configure and start the gRPC server.")
        parser.add_argument("--conf", required=True, help="Path to the service configuration file.")
        args = parser.parse_args()
        config_path = args.conf
        cls.serviceConfig = load_vapusaiserver_config(config_path)
        settings.set_service_config(cls.serviceConfig)
        service_logger.info("Loaded service config")
        #cls.secretsConfig = init_vapus_backend_secrets(cls.serviceConfig.mainConfig.vapusBEDbStore.path, cls.serviceConfig.mainConfig.vapusBESecretStore.path)
        #settings.set_secret_store(cls.secretsConfig)
        service_logger.info("Loaded secrets")
        cls.configure_logger(sys.argv)
        
    
    def start(cls):
        """
        Starts the gRPC server.
        """
        cls.configure()
        cls.init_server()
        egnine = cls.startDBEngine()
        # ServerBoot.boot(service_logger)

        try:
            cls.server.start()
            service_logger.info("Server started on port {port}", port=str(cls.serviceConfig.networkConfig.mlutilitySvc.port))
            cls.server.wait_for_termination()
        except Exception as e:
            service_logger.error("Error starting grpc server: {error}", error=str(e))
            raise e