import pytest

pytestmark = pytest.mark.django_db


class TestCategoryModel:
    def test_str_method(self, category_factory):
        # Arrange
        # Act
        x = category_factory(name="cate")
        # Assert
        assert x.__str__() == "cate"


class TestBrandModel:
    def test_str_method(self, brand_factory):
        x = brand_factory(name="brand")
        assert x.__str__() == "brand"


class TestProductModel:
    def test_str_method(self, product_factory):
        x = product_factory(name="prod")
        assert x.__str__() == "prod"
