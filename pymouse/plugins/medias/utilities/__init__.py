from hydrogram.types import CallbackQuery
from pymouse import PyMouse, YT_DLP, log

async def YouTubeDownloadAndSendMedia(
    c: PyMouse, # type: ignore
    cb: CallbackQuery,
    chat_id: int,
    vid: str,
    thumb: str,
    filepath: str,
    media_type: str,
    format_uid: str,
    i18n: dict
):
    # Download and GET infos of the video from YouTube
    yt = await YT_DLP().downloader_method(
        cb=cb,
        vid=vid,
        filepath=filepath,
        media_type=media_type,
        frt=format_uid,
        i18n=i18n
    )
    # Upload Media to Telegram using method (send_media) of the YT_DLP class
    await YT_DLP().send_media(
        c=c,
        cb=cb,
        chat_id=chat_id,
        thumb=thumb,
        media_type=media_type,
        media_class=yt,
        i18n=i18n,
    )