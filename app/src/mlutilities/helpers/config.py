from pydantic_settings import BaseSettings, SettingsConfigDict
from pathlib import Path
from dataclasses import dataclass
from utils.file_loader import load_basic_config
import os
class ServerConfig(BaseSettings):
    """
    Configuration class for server settings.
    """
    serviceName: str
    port: str
    scheme: str
    externalUrl: str
    svcType: str
    addr: str
    servicePort: str
    httpGwPort: str = "0000"
    
    class Config:
        case_sensitive = False

class ConfigBasePaths(BaseSettings):
    """
    Represents the configuration base paths.

    Attributes:
        path (str): The base path.

    Config:
        case_sensitive (bool): Whether the configuration is case sensitive or not.
    """
    filePath: str
    secret: str
    class Config:
        case_sensitive = False

class ServerCerts(BaseSettings):
  

  mtls: bool
  caCertFile: str
  serverCertFile: str
  serverKeyFile: str
  clientCertFile: str
  clientKeyFile: str

  class Config:
    case_sensitive = False
class ServiceConfig(BaseSettings):
    """
    Configuration class for the service.
    """

    vapusBESecretStorage: ConfigBasePaths
    vapusFileStorage: ConfigBasePaths
    # vapusBEDbStore: ConfigBasePaths
    # serverConfig: ServerConfig
    JWTAuthnSecrets: ConfigBasePaths
    artifactStore: ConfigBasePaths
    networkConfigFile: str
    serverCerts: ServerCerts
    class Config:
        case_sensitive = False

class NetworkConfig(BaseSettings):
    platformSvc: ServerConfig
    aistudioSvc: ServerConfig
    webappSvc: ServerConfig
    nabhikserver: ServerConfig
    mlutilitySvc: ServerConfig
@dataclass(frozen=True)
class VapusAiConfig:
    """
    Configuration class for Vapus AI.
    
    Attributes:
        mainConfig (ServiceConfig): The main configuration for Vapus AI.
    """
    mainConfig: ServiceConfig
    networkConfig: NetworkConfig
    configPath: str 

def load_vapusaiserver_config(configPath: str) -> VapusAiConfig:
    """
    Load the VapusAiServer configuration from the specified config file.

    Args:
        configPath (str): The path to the configuration file.

    Returns:
        VapusAiConfig: The loaded VapusAiConfig object.

    """
    configData = load_basic_config(os.path.join(configPath,"configs/mlutilities-config.yaml"))
    networkConfigPath = os.path.join(configPath,configData.get("networkConfigFile"))
    networkConfig = load_basic_config(networkConfigPath)
    return VapusAiConfig(mainConfig=ServiceConfig(**configData),networkConfig=NetworkConfig(**networkConfig),configPath=configPath)