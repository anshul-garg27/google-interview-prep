import asyncio
import sys
from dotenv import load_dotenv
from loguru import logger

from instagram.flows.schedule import refresh_creator_profiles_insta, refresh_non_creator_profiles_insta

flows = [refresh_creator_profiles_insta, refresh_non_creator_profiles_insta]


async def trigger_scheduled_flow(flow_name: str):
    for f in flows:
        if f.__name__ == flow_name:
            logger.info(f"Triggering flow {f.__name__}")
            await f()
            return logger.error("Flow not found. Exiting.")


"""
python trigger.py refresh_creator_profiles_insta
python trigger.py refresh_non_creator_profiles_insta
"""
if __name__ == "__main__":
    load_dotenv()
    if len(sys.argv) != 2:
        sys.exit("Invalid Arguments. Please specify flow name")
    flow = sys.argv[1]
    asyncio.run(trigger_scheduled_flow(flow))
