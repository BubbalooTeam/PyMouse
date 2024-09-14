#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.


import json
from time import sleep
from typing import Optional

class DataBase:
    def __init__(self):
        self.db_file = "pymouse/database/filepaths/database.json"
        self.load_db()

    def load_db(self):
        """
        Loader of DataBase File represented in `filepaths/database.json`.

        Returns:
            `dict`: Representing DataBase content.
        """
        try:
            with open(self.db_file, 'r') as file:
                self.db = json.load(file)
        except (FileNotFoundError):
            self.db = {}
        except (json.JSONDecodeError, UnicodeDecodeError):
            for try_again in range(100):
                sleep(1)
                print("There was an error decoding the DataBase... Trying again... in {range_again}/100 attempts.".format(range_again=try_again))
                try:
                    with open(self.db_file, 'r') as file:
                        self.db = json.load(file)
                    break
                except (FileNotFoundError):
                    self.db = {}
                except (json.JSONDecodeError, UnicodeDecodeError):
                    continue
            self.db = {}

    def save_db(self):
        """
        Save DataBase informations in `filepaths/database.json`.
        """
        try:
            with open(self.db_file, 'w') as file:
                json.dump(self.db, file, indent=4)
        except (FileNotFoundError):
            with open(self.db_file, "w+") as file:
                file.write(json.dumps(self.db, indent=4))

    class GetCollection:
        def __init__(self, collection: str):
            self.parent_db = DataBase()
            self.collection = collection
            if collection not in self.parent_db.db:
                self.parent_db.db[collection] = []

        def find(self) -> list:
            """
            Fetches a list of information in the current collection.

            Returns:
                `list`: Representing collection content.
            """
            return self.parent_db.db[self.collection]

        def find_one(self, filter: Optional[dict] = None) -> dict:
            """
            Fetches information in a collection according to the specified filter.

            Arguments:
                `filter (dict)`: Filter for Search informations.

            Returns:
                `dict`: Representing content of specified filter.
            """
            # find global collection info
            collection_data = self.find()
            if not filter:
                print("[database/modules]: No information provided for find especified info...")
                return
            # find keys
            try:
                index = next((i for i, item in enumerate(collection_data) if all(item[k] == v for k, v in filter.items())), None)
            except (IndexError, KeyError):
                index = None
            # find the information base with index
            if index is not None:
                return collection_data[index]
            else:
                return {}

        def insert_or_update(self, filter: Optional[dict] = None, info: Optional[dict] = None) -> bool:
            """
            Inserts or updates the information for the current collection.

            Arguments:
                `filter (dict)`: Filter to check for its existence in the DataBase.
                `info (dict)`: Information to be added/updated.
            """
            if not info:
                print("[database/modules]: No information provided for update...")
                return
            inup = False
            collection_data = self.find()
            if filter:
                # Find the index of the item to be updated
                try:
                    index = next((i for i, item in enumerate(collection_data) if all(item[k] == v for k, v in filter.items())), None)
                except (IndexError, KeyError):
                    index = None
                if index is not None:
                    # Update the existing item
                    collection_data[index].update(info)
                    inup = True
                else:
                    # Insert a new item
                    collection_data.append(info)
                    inup = True

            # Save the changes to the database
            self.parent_db.save_db()
            return inup

        def delete(self, filter: Optional[dict] = None) -> bool:
            """
            Remove a record from the database based on a specified key and its corresponding value.

            Arguments:
                `filter (dict)`: A dictionary containing the key-value pair for filtering records.

            Returns:
                `bool`: True if a record was deleted, False otherwise.
            """
            deleted = False
            if filter:
                finder = self.find()
                for record in finder:
                    if all(record.get(key) == value for key, value in filter.items()):
                        finder.remove(record)
                        deleted = True
            else:
                self.parent_db.db.clear()
                deleted = True
            # Save the changes to the database
            self.parent_db.save_db()
            return deleted