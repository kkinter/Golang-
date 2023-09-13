import pytest
from django.core.exceptions import ValidationError

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


class TestProductLineModel:
    def test_str_method(self, product_line_factory):
        obj = product_line_factory(sku="12345")
        assert obj.__str__() == "12345"

    def test_duplicate_order_values(self, product_line_factory, product_factory):
        obj = product_factory()
        product_line_factory(order=1, product=obj)
        with pytest.raises(ValidationError):
            product_line_factory(order=1, product=obj).clean()


class TestProductImageModel:
    def test_str_method(self, product_image_factory):
        obj = product_image_factory(order=1)
        assert obj.__str__() == "1"
