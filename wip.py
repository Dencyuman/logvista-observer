import logging
import time
import os
from pathlib import Path
from logging import getLogger, FileHandler, StreamHandler, Formatter
from logvista import LogvistaHandler

formatter = "%(levelname)-9s  %(asctime)s [%(filename)s:%(lineno)d] %(message)s"

logger = getLogger("auto_post")
logger.setLevel(logging.DEBUG)
sh = StreamHandler()
sh.setLevel(logging.DEBUG)
lv = LogvistaHandler("snapcheck-auto AAAAA")
lv.setLevel(logging.DEBUG)
logger.addHandler(sh)
logger.addHandler(lv)
sh.setFormatter(Formatter(formatter))

logger2 = getLogger("inspect")
logger2.setLevel(logging.DEBUG)
sh = StreamHandler()
sh.setLevel(logging.DEBUG)
lv = LogvistaHandler("snapcheck-scan/auto")
lv.setLevel(logging.DEBUG)
logger2.addHandler(sh)
logger2.addHandler(lv)
sh.setFormatter(Formatter(formatter))

try:
    raise Exception("test")
except Exception as e:
    logger.error("error")

logger.info("info")
logger2.info("info")
logger.warning("warning")
logger2.warning("warning")