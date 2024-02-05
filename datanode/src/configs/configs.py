import os
from decouple import Config, RepositoryEnv
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent
# ENV_NAME = os.environ['ENV'].lower()
DOTENV_FILE = os.path.abspath(f'{BASE_DIR}/test.env')
env_config = Config(RepositoryEnv(DOTENV_FILE))


class ConfigManager:
    configs = dict(
        server_port=env_config('DATANODE_PORT'),
        partition_count=env_config('PARTITIONS_COUNT'),
        partition_home_path=env_config('HOME_PATH'),
        encoding_method=env_config('ENCODING_METHOD'),
        pull_timeout=env_config('PULL_TIMEOUT'),
        pending_timeout=env_config('PENDING_TIMEOUT'),
        cleaner_period=env_config('CLEANER_PERIOD'),
        leader_host=env_config('LEADER_HOST'),
        leader_port=env_config('LEADER_PORT'),
        datanode_name=env_config('DATANODE_NAME'),
    )

    @staticmethod
    def get_prop(key):
        return ConfigManager.configs[key]
=======
        return ConfigManager.configs[key]
>>>>>>> origin/datanode-dev
