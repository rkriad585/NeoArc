from waitress import serve
from main import app
import config
import logging

logger = logging.getLogger("neoarc")

if __name__ == "__main__":
    logger.info("Starting NeoArc production server on %s:%s with waitress", config.HOST, config.PORT)
    serve(app, host=config.HOST, port=config.PORT)
