from pydantic import BaseSettings


class Settings(BaseSettings):
    heartbeat: bool = False

    class Config:
        env_file = ".env"
