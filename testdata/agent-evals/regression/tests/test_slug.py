import unittest

from slug import slugify


class SlugTests(unittest.TestCase):
    def test_words_are_lowercased_and_joined(self):
        self.assertEqual(slugify("Hello World"), "hello-world")

    def test_punctuation_is_removed_and_words_are_separated(self):
        self.assertEqual(slugify("Hello, World!"), "hello-world")


if __name__ == "__main__":
    unittest.main()
