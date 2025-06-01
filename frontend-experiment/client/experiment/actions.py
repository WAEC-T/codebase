from playwright.async_api import Page


async def press_link(page: Page, link_text: str):
    link_selector = f"a:text-is('{link_text}')"
    await page.click(link_selector)


async def submit(page: Page):
    await page.click(f"button[type='submit' i]")


async def submit_input(page: Page):
    await page.click(f"input[type='submit' i]")


async def press_button(page: Page, button_name: str):
    button_selector = f"button[name='{button_name}']"
    if not await page.is_visible(button_selector):
        button_selector = f"button:text-is('{button_name}')"

    if not button_selector or not await page.is_visible(button_selector):
        raise Exception(f"Button '{button_name}' not found")

    await page.click(button_selector)


async def navigate_to(page: Page, url: str):
    await page.goto(url)


async def fill_input(page: Page, input_name: str, value: str):
    input_selector = f"input[name='{input_name}' i]"
    await page.fill(input_selector, value)


async def get_text(page: Page, selector: str) -> str:
    return await page.text_content(selector)
