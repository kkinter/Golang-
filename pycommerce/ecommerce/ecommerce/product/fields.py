from django.core import checks
from django.core.exceptions import ObjectDoesNotExist
from django.db import models


class OrderField(models.PositiveIntegerField):
    description = "Ordering field on a unique field"

    def __init__(self, unique_for_field=None, *args, **kwargs):
        self.unique_for_field = unique_for_field
        super().__init__(*args, **kwargs)

    def check(self, **kwargs):
        return [
            *super().check(**kwargs),
            *self._check_for_field_attribute(**kwargs),
        ]

    def _check_for_field_attribute(self, **kwargs):
        if self.unique_for_field is None:
            return [checks.Error("OrderField는 'unique_for_field' 속성을 정의해야 합니다.")]
        elif self.unique_for_field not in [f.name for f in self.model._meta.get_fields()]:
            return [checks.Error("입력한 OrderField가 기존 모델 필드와 일치하지 않습니다.")]

        return []

    def pre_save(self, model_instance, add):
        print("HELLO")
        print(model_instance)

        if getattr(model_instance, self.attname) is None:
            # print(getattr(model_instance, self.attname))
            # print("NEED A VAL")
            qs = self.model.objects.all()
            try:
                query = {
                    self.unique_for_field: getattr(model_instance, self.unique_for_field),
                }
                # print(query)
                # {'product': <Product: p1>}
                qs = qs.filter(**query)
                # print(qs)
                # <ActiveQueryset [<ProductLine: 1>, <ProductLine: 1>]>
                last_item = qs.latest(self.attname)
                value = last_item.order + 1
                # print(self.attname)
                # order
            except ObjectDoesNotExist:
                value = 1

            return value
        else:
            return super().pre_save(model_instance, add)
