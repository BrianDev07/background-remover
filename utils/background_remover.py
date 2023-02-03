from PIL import Image
from tkinter import Tk
from tkinter.filedialog import askopenfilename
from os import path
from os.path import isfile


def __show_dialog() -> str:
    """Creates a file selection window and gets the path of the selected file.

    Returns:
        str: Absolute path of the selected file.
    """

    Tk().withdraw()  # keeps root window from appearing
    return askopenfilename()


def __trim_extension(file_path: str) -> str:
    """Takes the path of a file and returns its name with no extension.
    This is done so the newly created image has the same name, but not necesarily the same extension.
    i.e. C:\THESE\ARE\FOLDERS\My_file.extension -> My_file

    Args:
        file_path (str): Absolute path of the file.

    Returns:
        str: Name of the file without extension.
    """

    name: str = file_path.split('.')
    filename: str = name[0].split('/')

    return filename[-1]    # last element


def to_black_and_white(threshold: int) -> tuple[Image.Image, str]:
    """Turns an image to grayscale. A threshold is used to control the black and white effect.
    Values above the threshold are turned white, and those below it are turned black.

    Args:
        threshold (int): Works similar to a brightness sensitivity regulator.

    Returns:
        Image.Image: Grayscale image.
        str: Absolute path where the file is to be saved.
    """

    file_path: str = __show_dialog()
    rgb_img: Image.Image = Image.open(file_path)
    gray_img: Image.Image = rgb_img.convert('L')         # 'L' is for Luminance. Converts image to grayscale

    return gray_img.point(lambda x: 0 if x < threshold else 255), file_path


def remove_background(image, file_path) -> tuple[bool, str]:
    """Takes a black and white image and removes its background.
    If the a pixel is white, meaning its RGBA values are [255, 255, 255, alpha], alpha is set to 0,
    then if alpha is 0, it means that the color is transparent. This way, all white is removed.

    Args:
        image (Image.Image): Image that is being processed.
        file_path (str): Absolute path where the image is stored.

    Returns:
        bool: Indicates if the edited image was created.
        str: Absolute path where the edited image was created.
    """

    current_dir: str = path.dirname(file_path)
    image: Image.Image = image.convert("RGBA")  # alpha (A) specifies the opacity of a color
    pixel_collection = image.getdata()          # stores pixel information of every pixel
    new_data: list = []                         # list that will store the pixel data of the new image

    for pixel in pixel_collection:
        if pixel[0] == 255 and pixel[1] == 255 and pixel[2] == 255:
            new_data.append((255, 255, 255, 0)) # alpha is set to 0 here
        else:
            new_data.append(pixel)  # if the pixel is now white, appends it as it is (black).

    image.putdata(new_data) # copies pixel data from the data sequence into the image

    filename_no_extension: str = __trim_extension(file_path)
    dest_filepath: str = f"{current_dir}/{filename_no_extension}-no_background.png"

    image.save(dest_filepath, None)  # creates an image with no background

    return isfile(dest_filepath), dest_filepath
