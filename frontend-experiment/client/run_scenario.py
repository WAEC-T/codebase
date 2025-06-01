import asyncio
from experiment.scenario import execute_scenario, makeNameList, run, setup_logging
import sys

async def main():
    if len(sys.argv) < 2:
        raise ValueError("Usage: python scenario.py <url> <headless> <slow_mo>")

    url = sys.argv[1]
    headless = sys.argv[2].lower() == "true" if len(sys.argv) > 2 else True
    slow_mo = int(sys.argv[3]) if len(sys.argv) > 3 else None

    nameList = makeNameList()

    # Make logfile
    setup_logging()

    await execute_scenario(nameList, url=url, headless=headless, slow_mo=slow_mo)

if __name__ == "__main__":
    asyncio.run(main())