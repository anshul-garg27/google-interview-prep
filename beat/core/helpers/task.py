import traceback

from loguru import logger



def success(name, *args, **kwargs):
    logger.error(f"{name} {args} {kwargs}")


def error(name, err, *args, **kwargs):
    logger.error(f"{name} {args} {kwargs}")


def emit(func):
    def wrapper(*args, **kwargs):
        try:
            func_res = func(*args, **kwargs)
            success(func.__name__, *args, **kwargs)
        except Exception as e:
            err = traceback.format_exc()
            error(func.__name__, err, *args, **kwargs)
            raise e
        finally:
            return func_res
    return wrapper
