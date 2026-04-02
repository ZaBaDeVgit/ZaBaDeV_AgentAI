---
name: pytest
description: >
  Pytest testing patterns for Python. Trigger: When writing Python tests - fixtures, mocking, markers.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- Writing unit tests for Python code
- Creating fixtures for test setup
- Mocking external dependencies
- Parameterizing tests with multiple inputs

---

## Critical Patterns

### Pattern 1: Fixtures

```python
import pytest
from myapp import create_app, db

# Function-scoped fixture
@pytest.fixture
def app():
    app = create_app('testing')
    with app.app_context():
        db.create_all()
        yield app
        db.drop_all()

# Session-scoped fixture (runs once)
@pytest.fixture(scope='session')
def db_engine():
    engine = create_engine('sqlite:///:memory:')
    yield engine
    engine.dispose()

# Fixture with dependencies
@pytest.fixture
def client(app):
    return app.test_client()
```

### Pattern 2: Fixture with teardown

```python
@pytest.fixture
def temp_dir(tmp_path):
    """Create temp directory, cleaned up after."""
    test_dir = tmp_path / 'test_data'
    test_dir.mkdir()
    
    # Write some test files
    (test_dir / 'file.txt').write_text('test content')
    
    yield test_dir
    # Cleanup automatic with tmp_path
```

### Pattern 3: Mocking

```python
from unittest.mock import patch, MagicMock

# Patch external calls
@patch('myapp.external_api.call')
def test_external_api(mock_call):
    mock_call.return_value = {'status': 'ok'}
    
    result = myapp.process()
    
    mock_call.assert_called_once()

# Mock class
class TestUserService:
    @patch('myapp.services.UserService')
    def test_create_user(self, mock_service):
        mock_instance = MagicMock()
        mock_instance.create.return_value = {'id': 1, 'name': 'Test'}
        mock_service.return_value = mock_instance
        
        result = create_user('Test')
        
        assert result['name'] == 'Test'
```

### Pattern 4: pytest-mock

```python
def test_with_mocker(mocker):
    # Mock via mocker fixture
    mock = mocker.patch('myapp.send_email')
    
    result = register_user('test@test.com')
    
    mock.assert_called_once()

# Spy on existing methods
def test_spy(mocker):
    spy = mocker.spy(myapp.utils, 'format_date')
    
    result = process_date('2024-01-01')
    
    spy.assert_called_once()
```

### Pattern 5: Markers

```python
import pytest

# Skip test
@pytest.mark.skip(reason='Not implemented yet')
def test_future_feature():
    pass

# Skip if condition
@pytest.mark.skipif(sys.version_info < (3, 10), reason='Requires Python 3.10+')
def test_new_feature():
    pass

# Expected failure
@pytest.mark.xfail(reason='Known issue #123')
def test_known_bug():
    assert False

# Parameterize
@pytest.mark.parametrize('input,expected', [
    ('hello', 'HELLO'),
    ('world', 'WORLD'),
    ('test', 'TEST'),
])
def test_uppercase(input, expected):
    assert input.upper() == expected
```

### Pattern 6: Async Testing

```python
import pytest
import asyncio

# pytest-asyncio mode='auto'
@pytest.mark.asyncio
async def test_async_function():
    result = await async_fetch_data()
    assert result is not None

@pytest.mark.asyncio
async def test_with_event_loop():
    async with aiohttp.ClientSession() as session:
        response = await session.get('http://test.com')
        assert response.status == 200
```

---

## Decision Tree

```
Test setup needed?
├── Simple → No fixture
├── Reusable → @pytest.fixture
├── Expensive → @pytest.fixture(scope='session')
└── Complex → Factory fixture

Mocking?
├── External API → @patch
├── Class → MagicMock
├── Existing method → mocker.spy
└── pytest-mock → mocker fixture

Execution control?
├── Skip → @pytest.mark.skip
├── Skip condition → @pytest.mark.skipif
├── Expected fail → @pytest.mark.xfail
└── Multiple inputs → @pytest.mark.parametrize
```

---

## Anti-Patterns

- ❌ Not using fixtures - leads to duplicated setup
- ❌ Over-mocking - test the real thing when possible
- ❌ Using assert in setup code - use fixtures
- ❌ Not using parametrize - duplicate test functions

---

## Commands

```bash
pytest                          # Run all tests
pytest -v                       # Verbose
pytest -k "test_name"          # Run matching tests
pytest --collect-only          # List tests
pytest --markers               # List markers
pytest --fixture fixtures.py   # List fixtures
pytest -x                      # Stop on first failure
pytest --lf                    # Run last failed only
```

---

## Resources

- **Pytest Docs**: https://docs.pytest.org/
- **Fixtures**: https://docs.pytest.org/en/latest/fixture.html
- **pytest-mock**: https://pytest-mock.readthedocs.io/