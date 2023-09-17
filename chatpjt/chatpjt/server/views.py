from django.db.models import Count
from rest_framework import viewsets
from rest_framework.exceptions import AuthenticationFailed, ValidationError
from rest_framework.response import Response

from .models import Server
from .schema import server_list_docs
from .serializer import ServerSerializer


class ServerListViewSet(viewsets.ViewSet):
    queryset = Server.objects.all()

    @server_list_docs
    def list(self, request):
        """다양한 매개변수로 필터링된 서버 목록을 반환합니다.

        이 메서드는 request 객체에 제공된 쿼리 매개변수를 기반으로 서버의 쿼리셋을 검색합니다.
        다음과 같은 쿼리 매개변수가 지원됩니다:

        category: 카테고리 이름으로 서버를 필터링합니다.
        qty: 반환되는 서버의 수를 제한합니다.
        by_user: 사용자 ID에 따라 서버를 필터링하며, 사용자가 멤버인 서버만 반환합니다.
        by_serverid: 서버 ID로 서버를 필터링합니다.
        with_num_members: 각 서버에 멤버 수를 주석으로 추가합니다.
        인수:
        request: 쿼리 매개변수가 포함된 Django Request 객체입니다.

        반환값:
        지정된 매개변수로 필터링된 서버의 쿼리셋입니다.

        예외:
        AuthenticationFailed: 'by_user' 또는 'by_serverid' 매개변수가 쿼리에 포함되어 있고 사용자가 인증되지 않은 경우 발생합니다.
        ValidationError: 쿼리 매개변수의 구문 분석 또는 유효성 검사 오류가 있는 경우 발생합니다.
        'by_serverid' 매개변수가 유효한 정수가 아니거나 지정된 ID의 서버가 없는 경우 발생할 수 있습니다.

        예시:
        '게임' 카테고리에서 적어도 5명의 멤버를 가진 모든 서버를 검색하려면 다음 요청을 수행할 수 있습니다:

        GET /servers/?category=gaming&with_num_members=true&num_members__gte=5

        인증된 사용자가 멤버인 첫 10개의 서버를 검색하려면 다음 요청을 수행할 수 있습니다:

        GET /servers/?by_user=true&qty=10

        """
        category = request.query_params.get("category")
        qty = request.query_params.get("qty")
        by_user = request.query_params.get("by_user") == "true"
        by_serverid = request.query_params.get("by_serverid")
        with_num_members = request.query_params.get("with_num_members") == "true"

        if category:
            self.queryset = self.queryset.filter(category__name=category)

        if by_user:
            if by_user and request.user.is_authenticated:
                user_id = request.user.id
                self.queryset = self.queryset.filter(member=user_id)
            else:
                raise AuthenticationFailed()

        if with_num_members:
            self.queryset = self.queryset.annotate(num_members=Count("member"))

        if by_serverid:
            if not request.user.is_authenticated:
                raise AuthenticationFailed()

            try:
                self.queryset = self.queryset.filter(id=by_serverid)
                if not self.queryset.exists():
                    raise ValidationError(detail=f"Server with id {by_serverid} not found")
            except ValueError:
                raise ValidationError(detail="Server value error")

        if qty:
            self.queryset = self.queryset[: int(qty)]

        serializer = ServerSerializer(self.queryset, many=True, context={"num_members": with_num_members})
        return Response(serializer.data)
