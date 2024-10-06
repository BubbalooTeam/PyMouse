class FMT:
    @staticmethod
    def FormatListtoString(RelativeList: list, i18n: dict) -> str:
        """
        Convert Lists to string format.

        Arguments:
            - `RelativeList(list)`: Relative list with several strings, which must be transformed into string.
            - `i18n(dict)`: Dictionary of international words.

        Returns:
            - `ListtoString(str)`: List converted into string.
        """
        ListtoString = (
            ", ".join(RelativeList[:-1])
            + i18n["fmt"]["and"] + RelativeList[-1]
            if len(RelativeList) > 1
            else RelativeList[0]
            if RelativeList else ""
        )
        return ListtoString