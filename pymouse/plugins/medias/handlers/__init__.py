from shutil import rmtree
from os import path
from re import compile
from wget import download
from tempfile import TemporaryDirectory
from youtubesearchpython.__future__ import VideosSearch
from traceback import format_exc

from hydrogram.types import Message, CallbackQuery, InputMediaPhoto, InputMediaAudio, InputMediaVideo
from hydrogram.errors import MessageNotModified
from hydrogram.enums import ChatAction
from hydrokeyboard import InlineKeyboard, InlineButton

from pymouse import PyMouse, Decorators, Config, log, YT_DLP, DownloadPaths
from pymouse.plugins.medias.utilities import YouTubeDownloadAndSendMedia
from pymouse.utils import HandleText

YOUTUBE_REGEX = compile(
    r"(?m)http(?:s?):\/\/(?:www\.)?(?:music\.)?youtu(?:be\.com\/(watch\?v=|shorts/|embed/)|\.be\/|)(?P<id>([\w\-\_]{11}))(&(amp;)?‌​[\w\?‌​=]*)?"
)

YT_VAR = {}

class Medias_Plugins:
    @staticmethod
    @Decorators().Locale()
    async def ytdl_handler(_, m: Message, i18n): # type: ignore
        query = HandleText().input_str(union=m)
        chat_id = m.chat.id
        keyboard = InlineKeyboard()
        # //
        if not query:
            return await m.reply(i18n["youtube-dl"]["checkers"]["noArgs"])
        match = YOUTUBE_REGEX.match(query)
        if match is None:
            search_key = await YT_DLP().rand_key()
            YT_VAR[search_key] = query
            search = await VideosSearch(query).next()
            if search["result"] == []:
                return await m.reply(f"No result found for `{query}`")
            i = search["result"][0]
            out = f"<b><a href={i['link']}>{i['title']}</a></b>"
            out += f"\nPublished {i['publishedTime']}\n"
            out += f"\n<b>❯ Duration:</b> {i['duration']}"
            out += f"\n<b>❯ Views:</b> {i['viewCount']['short']}"
            out += f"\n<b>❯ Uploader:</b> <a href={i['channel']['link']}>{i['channel']['name']}</a>\n\n"
            # //
            keyboard.paginate(len(search["result"])-1, 1, f"ytdl_scroll|{search_key}" + "|{number}|" + f"{m.from_user.id}")
            keyboard.row(InlineButton("Download", f"yt_gen|{i['id']}|{None}|{m.from_user.id}"))
            # //
            img = await YT_DLP().get_ytthumb(i["id"])
            caption = out
            markup = keyboard
            await m.reply_photo(photo=img, caption=caption, reply_markup=markup)
        else:
            key = match.group("id")
            deatils = await YT_DLP().get_download_button(key, m.from_user.id)
            img = await YT_DLP().get_ytthumb(key)
            caption = deatils.caption
            markup = deatils.buttons
            await m.reply_photo(photo=img, caption=caption, reply_markup=markup)

    @staticmethod
    @Decorators().Locale()
    async def ytdl_scroll_callback(c: PyMouse, cb: CallbackQuery, i18n): # type: ignore
        inf = cb.data.split("|")
        search_key = inf[1]
        page = int(inf[2])
        user_id = int(inf[3])
        keyboard = InlineKeyboard()
        # //
        if not cb.from_user.id == user_id:
            return await cb.answer(i18n["youtube-dl"]["checkers"]["notforYou"], show_alert=True)
        try:
            query = YT_VAR[search_key]
            search = await VideosSearch(query).next()
            i = search["result"][page]
            out = f"<b><a href={i['link']}>{i['title']}</a></b>"
            out += f"\nPublished {i['publishedTime']}\n"
            out += f"\n<b>❯ Duration:</b> {i['duration']}"
            out += f"\n<b>❯ Views:</b> {i['viewCount']['short']}"
            out += f"\n<b>❯ Uploader:</b> <a href={i['channel']['link']}>{i['channel']['name']}</a>\n\n"
            # //
            keyboard.paginate(len(search["result"])-1, page, f"ytdl_scroll|{search_key}" + "|{number}|" + f"{user_id}")
            keyboard.row(InlineButton("Download", f"yt_gen|{i['id']}|{None}|{user_id}"))

            if page == 0:
                if len(search["result"]) == 1:
                    return await cb.answer("That's the end of list", show_alert=True)

            await cb.edit_message_media(
                InputMediaPhoto(await YT_DLP().get_ytthumb(i["id"]), caption=out), reply_markup=keyboard
            )
        except (IndexError, KeyError, MessageNotModified):
            return await cb.answer(
                "error when obtaining information, perform a new search", show_alert=True
            )

    @staticmethod
    @Decorators().Locale()
    async def download_handler(c: PyMouse, cb: CallbackQuery, i18n): # type: ignore
        inf = cb.data.split("|")
        key = inf[1]
        user_id = int(inf[3])
        # //
        try:
            chat_id = cb.message.chat.id
        except Exception:
            chat_id = int(inf[3])
        if not cb.from_user.id == user_id:
            return await cb.answer(i18n["youtube-dl"]["checkers"]["notforYou"], show_alert=True)
        try:
            if inf[0] == "yt_gen":
                x = (await YT_DLP().get_download_button(
                    yt_id=key,
                    user_id=user_id
                    )
                )
                await cb.edit_message_caption(caption=x.caption, reply_markup=x.buttons)
            else:
                uid = inf[2]
                type_ = inf[4]
                with TemporaryDirectory(dir=Config.DOWNLOAD_PATH) as tempdir:
                    path_ = path.join(tempdir, "ytdl")
                thumb = download(await YT_DLP().get_ytthumb(key), Config.DOWNLOAD_PATH + "thumbnail.png")
                if type_ == "a":
                    format_ = "audio"
                else:
                    format_ = "video"

                await YouTubeDownloadAndSendMedia(
                    c=c,
                    cb=cb,
                    chat_id=chat_id,
                    vid=key,
                    thumb=thumb,
                    filepath=path_,
                    media_type=format_,
                    format_uid=uid,
                    i18n=i18n,
                )

                await c.delete_messages(
                    chat_id=chat_id,
                    message_ids=cb.message.id,
                )
                DownloadPaths().Delete_FileName(filename=thumb)
                rmtree(tempdir)
        except MessageNotModified:
            return
        except Exception as e:
            err = format_exc(limit=3020)
            log.error(err)
            return