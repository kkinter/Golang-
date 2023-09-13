from drf_spectacular.utils import extend_schema
from rest_framework import viewsets
from rest_framework.decorators import action
from rest_framework.response import Response

from .models import Brand, Category, Product
from .serializers import BrandSerializer, CategorySerializer, ProductSerializer


class CategoryViewSet(viewsets.ViewSet):
    """
    모든 카테고리 조회를 위한 Viewset
    """

    queryset = Category.objects.all()

    @extend_schema(responses=CategorySerializer)
    def list(self, request):
        serializer = CategorySerializer(self.queryset, many=True)
        return Response(serializer.data)


class BrandViewSet(viewsets.ViewSet):
    """
    모든 브랜드 조회를 위한 Viewset
    """

    queryset = Brand.objects.all()

    @extend_schema(responses=BrandSerializer)
    def list(self, request):
        serializer = BrandSerializer(self.queryset, many=True)
        return Response(serializer.data)


class ProductViewSet(viewsets.ViewSet):
    """
    모든 상품 조회를 위한 Viewset
    """

    queryset = Product.objects.all()
    lookup_field = "slug"

    def retrieve(self, request, slug=None):
        """
        단일 상품 조회를 위한 Endpoint
        """
        serializer = ProductSerializer(self.queryset.filter(slug=slug), many=True)
        return Response(serializer.data)

    @extend_schema(responses=ProductSerializer)
    def list(self, request):
        serializer = ProductSerializer(self.queryset, many=True)
        return Response(serializer.data)

    @action(
        methods=["get"],
        detail=False,
        url_path=r"category/(?P<category>\w+)/all",
    )
    def list_product_by_category(self, request, category=None):
        """
        카테고리 별 상품 조회를 위한 Endpoint
        """
        serializer = ProductSerializer(self.queryset.filter(category__name=category), many=True)

        return Response(serializer.data)
