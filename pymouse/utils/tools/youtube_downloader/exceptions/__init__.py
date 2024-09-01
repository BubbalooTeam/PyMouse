class YouTubeDownloaderError(Exception):
    """Raise if any error occurs in YouTube Downloader"""

class YouTubeExtractError(YouTubeDownloaderError):
    """Raise if an error occurs while extracting information from the video."""