from PIL import UnidentifiedImageError
from utils import ( to_black_and_white, remove_background )


def main() -> None:
    try:
        threshold: int = 128
        
        (black_and_white, file_path) = to_black_and_white(threshold)
        (file_exists, dest_path) = remove_background(black_and_white, file_path)
    except (TypeError, AttributeError, UnboundLocalError, UnidentifiedImageError):
        print("ERROR: NO IMAGE WAS PROVIDED OR SELECTED FILE WAS NOT AN IMAGE.")
        return
    
    print(f"IMAGE CREATED IN {dest_path}" if file_exists else "ERROR: IMAGE WAS NOT CREATED.")


if __name__ == '__main__':
    main()
