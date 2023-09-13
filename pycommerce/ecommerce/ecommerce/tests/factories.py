import factory

from ecommerce.product.models import Brand, Category, Product


class CategoryFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = Category

    # name = "test_category"
    name = factory.Sequence(lambda n: "Category_%d" % n)


class BrandFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = Brand

    # name = "test_brand"
    name = factory.Sequence(lambda n: "Brand_%d" % n)


class ProductFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = Product

    # name = "test_product"
    name = factory.Sequence(lambda n: "Product_%d" % n)
    description = "test_description"
    is_digital = True
    # 외래키 참조
    brand = factory.SubFactory(BrandFactory)
    category = factory.SubFactory(CategoryFactory)
