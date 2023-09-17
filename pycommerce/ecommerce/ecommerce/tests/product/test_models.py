import pytest
from django.core.exceptions import ValidationError

from ecommerce.product.models import ProductTypeAttribute

pytestmark = pytest.mark.django_db


class TestCategoryModel:
    def test_str_method(self, category_factory):
        x = category_factory(name="cate")
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
    def test_str_method(self, product_line_factory, attribute_value_factory):
        attr = attribute_value_factory(attribute_value="test")
        obj = product_line_factory.create(sku="12345", attribute_value=(attr,))
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


class TestProductTypeModel:
    def test_str_method(self, product_type_factory, attribute_factory):
        test = attribute_factory(name="test")
        obj = product_type_factory.create(name="test_type", attribute=(test,))

        x = ProductTypeAttribute.objects.get(id=1)
        print(x)

        assert obj.__str__() == "test_type"


class TestAttributeModel:
    def test_str_method(self, attribute_factory):
        obj = attribute_factory(name="test_attribute")
        assert obj.__str__() == "test_attribute"


class TestAttributeValueModel:
    def test_str_method(self, attribute_value_factory, attribute_factory):
        obj_a = attribute_factory(name="test_attribute")
        obj_b = attribute_value_factory(attribute_value="test_value", attribute=obj_a)
        assert obj_b.__str__() == "test_attribute-test_value"
