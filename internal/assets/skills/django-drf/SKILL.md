---
name: django-drf
description: >
  Django REST Framework patterns. Trigger: When building REST APIs with Django - ViewSets, Serializers, Filters.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- Building REST APIs with Django
- Creating serializers for data validation
- Implementing authentication
- Adding filtering and pagination

---

## Critical Patterns

### Pattern 1: Serializers

```python
from rest_framework import serializers

# ModelSerializer - auto-generates from model
class UserSerializer(serializers.ModelSerializer):
    class Meta:
        model = User
        fields = ['id', 'username', 'email', 'date_joined']
        read_only_fields = ['id', 'date_joined']

# Custom serializer with validation
class ArticleSerializer(serializers.ModelSerializer):
    author = serializers.PrimaryKeyRelatedField(read_only=True)
    
    class Meta:
        model = Article
        fields = ['id', 'title', 'content', 'author', 'created_at']
    
    def validate_title(self, value):
        if len(value) < 10:
            raise serializers.ValidationError("Title too short")
        return value
    
    def validate(self, data):
        if data.get('published') and not data.get('content'):
            raise serializers.ValidationError("Published articles need content")
        return data
```

### Pattern 2: ViewSets and Routers

```python
from rest_framework import viewsets, routers

# ModelViewSet - full CRUD
class UserViewSet(viewsets.ModelViewSet):
    queryset = User.objects.all()
    serializer_class = UserSerializer
    permission_classes = [IsAuthenticated]
    
    def get_queryset(self):
        # Filter by current user
        return User.objects.filter(id=self.request.user.id)

# ReadOnlyViewSet
class ArticleViewSet(viewsets.ReadOnlyModelViewSet):
    queryset = Article.objects.all()
    serializer_class = ArticleSerializer

# Router registration
router = routers.DefaultRouter()
router.register(r'users', UserViewSet)
router.register(r'articles', ArticleViewSet)

urlpatterns = [
    path('', include(router.urls)),
]
```

### Pattern 3: Custom Actions

```python
class OrderViewSet(viewsets.ModelViewSet):
    queryset = Order.objects.all()
    
    @action(detail=True, methods=['post'])
    def cancel(self, request, pk=None):
        order = self.get_object()
        order.status = 'cancelled'
        order.save()
        return Response({'status': 'cancelled'})
    
    @action(detail=False, methods=['get'])
    def summary(self, request):
        """Get order summary for current user."""
        orders = self.get_queryset()
        return Response({
            'total': orders.count(),
            'total_value': sum(o.total for o in orders),
        })
```

### Pattern 4: Filters

```python
from rest_framework import filters

# Filter backends
class ArticleViewSet(viewsets.ModelViewSet):
    queryset = Article.objects.all()
    filter_backends = [
        filters.SearchFilter,
        filters.OrderingFilter,
        DjangoFilterBackend,
    ]
    filterset_class = ArticleFilter
    search_fields = ['title', 'content']
    ordering_fields = ['created_at', 'published']
    ordering = ['-created_at']

# Custom filter
from django_filters import FilterSet, rest_framework as df_filters

class ArticleFilter(FilterSet):
    author = df_filters.NumberFilter(field_name='author__id')
    published = df_filters.BooleanFilter(field_name='published')
    category = df_filters.ChoiceFilter(choices=Category.choices)
    
    class Meta:
        model = Article
        fields = ['author', 'published', 'category']
```

### Pattern 5: Authentication

```python
from rest_framework import authentication, permissions

# Token authentication
from rest_framework.authtoken.views import ObtainAuthToken
from rest_framework.authtoken.models import Token

class CustomAuthToken(ObtainAuthToken):
    def post(self, request):
        serializer = self.serializer_class(
            data=request.data,
            context={'request': request}
        )
        serializer.is_valid(raise_exception=True)
        user = serializer.validated_data['user']
        token, created = Token.objects.get_or_create(user=user)
        return Response({'token': token.key})

# JWT Authentication (djangorestframework-simplejwt)
from rest_framework_simplejwt.views import TokenObtainPairView, TokenRefreshView

urlpatterns = [
    path('api/token/', TokenObtainPairView.as_view(), name='token_obtain_pair'),
    path('api/token/refresh/', TokenRefreshView.as_view(), name='token_refresh'),
]

# Custom permission
class IsOwnerOrReadOnly(permissions.BasePermission):
    def has_object_permission(self, request, view, obj):
        if request.method in permissions.SAFE_METHODS:
            return True
        return obj.owner == request.user
```

### Pattern 6: Nested Serializers

```python
class AuthorSerializer(serializers.ModelSerializer):
    class Meta:
        model = Author
        fields = ['id', 'name']

class BookSerializer(serializers.ModelSerializer):
    author = AuthorSerializer(read_only=True)
    author_id = serializers.PrimaryKeyRelatedField(
        queryset=Author.objects.all(),
        source='author',
        write_only=True
    )
    
    class Meta:
        model = Book
        fields = ['id', 'title', 'author', 'author_id']
```

---

## Decision Tree

```
API structure?
├── CRUD → ModelViewSet
├── Read-only → ReadOnlyModelViewSet
├── Custom actions → @action decorator
└── Multiple resources → Nested routers

Authentication?
├── Simple → TokenAuthentication
├── JWT → djangorestframework-simplejwt
└── OAuth → django-oauth-toolkit

Filtering?
├── Simple → SearchFilter, OrderingFilter
├── Complex → django-filter
└── Custom → FilterSet class
```

---

## Anti-Patterns

- ❌ Not using ModelSerializer - loses validation
- ❌ Putting business logic in serializers
- ❌ Not adding pagination - causes performance issues
- ❌ Overriding get_queryset without considering filters

---

## Commands

```bash
python manage.py startapp api          # Create API app
python manage.py makemigrations         # Create migrations
python manage.py migrate                # Run migrations
python manage.py shell                  # Django shell
```

---

## Resources

- **DRF Docs**: https://www.django-rest-framework.org/
- **Django Filter**: https://django-filter.readthedocs.io/
- **SimpleJWT**: https://django-rest-framework-simplejwt.readthedocs.io/