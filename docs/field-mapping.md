# Карта полей анкеты → Шаблоны банка

## Обзор

Данный документ описывает соответствие полей анкеты заявки с полями, требуемыми различными банками для подачи заявок на кредитование и банковские гарантии.

## Структура анкеты

### 1. Личные данные (Personal Data)

| Поле анкеты | Описание | Банк 1 | Банк 2 | Банк 3 | Обязательность |
|-------------|----------|--------|--------|--------|----------------|
| `lastName` | Фамилия | `surname` | `last_name` | `family_name` | Обязательно |
| `firstName` | Имя | `name` | `first_name` | `given_name` | Обязательно |
| `middleName` | Отчество | `patronymic` | `middle_name` | `patronymic` | Опционально |
| `birthDate` | Дата рождения | `birth_date` | `date_of_birth` | `birthday` | Обязательно |
| `birthPlace` | Место рождения | `birth_place` | `place_of_birth` | `birth_location` | Обязательно |
| `gender` | Пол | `gender` | `sex` | `gender` | Обязательно |
| `citizenship` | Гражданство | `citizenship` | `nationality` | `citizenship` | Обязательно |
| `maritalStatus` | Семейное положение | `marital_status` | `family_status` | `marital_state` | Обязательно |
| `hasChildren` | Наличие детей | `has_children` | `children_count` | `dependents` | Обязательно |
| `childrenCount` | Количество детей | `children_count` | `number_of_children` | `children_count` | Опционально |
| `passportSeries` | Серия паспорта | `passport_series` | `passport_series` | `passport_series` | Обязательно |
| `passportNumber` | Номер паспорта | `passport_number` | `passport_number` | `passport_number` | Обязательно |
| `passportIssuedBy` | Кем выдан | `passport_issued_by` | `issuing_authority` | `passport_authority` | Обязательно |
| `passportIssueDate` | Дата выдачи | `passport_issue_date` | `issue_date` | `passport_date` | Обязательно |
| `passportDepartmentCode` | Код подразделения | `department_code` | `issuing_code` | `dept_code` | Обязательно |

### 2. Контактные данные (Contact Data)

| Поле анкеты | Описание | Банк 1 | Банк 2 | Банк 3 | Обязательность |
|-------------|----------|--------|--------|--------|----------------|
| `primaryPhone` | Основной телефон | `phone` | `mobile_phone` | `contact_phone` | Обязательно |
| `additionalPhones` | Доп. телефоны | `additional_phones[]` | `other_phones[]` | `alt_phones[]` | Опционально |
| `email` | Email | `email` | `email_address` | `email` | Обязательно |
| `registrationAddress` | Адрес регистрации | `reg_address` | `registered_address` | `legal_address` | Обязательно |
| `actualAddress` | Фактический адрес | `actual_address` | `residential_address` | `home_address` | Обязательно |
| `sameAddress` | Адреса совпадают | `addresses_match` | `same_residence` | `address_equal` | Опционально |

### 3. Профессиональные данные (Professional Data)

| Поле анкеты | Описание | Банк 1 | Банк 2 | Банк 3 | Обязательность |
|-------------|----------|--------|--------|--------|----------------|
| `currentJob.companyName` | Название компании | `employer` | `company_name` | `employer_name` | Обязательно |
| `currentJob.position` | Должность | `position` | `job_title` | `position` | Обязательно |
| `currentJob.workPhone` | Рабочий телефон | `work_phone` | `office_phone` | `business_phone` | Обязательно |
| `currentJob.workEmail` | Рабочий email | `work_email` | `business_email` | `work_email` | Опционально |
| `currentJob.workAddress` | Адрес работы | `work_address` | `office_address` | `business_address` | Обязательно |
| `currentJob.employmentDate` | Дата трудоустройства | `employment_date` | `start_date` | `hire_date` | Обязательно |
| `currentJob.monthlyIncome` | Месячный доход | `monthly_income` | `salary` | `income` | Обязательно |
| `education.level` | Уровень образования | `education_level` | `degree` | `education` | Обязательно |
| `education.institution` | Учебное заведение | `university` | `institution` | `school` | Обязательно |
| `education.graduationYear` | Год окончания | `graduation_year` | `completion_year` | `grad_year` | Обязательно |
| `education.specialty` | Специальность | `specialty` | `field_of_study` | `major` | Обязательно |

### 4. Финансовые данные (Financial Data)

| Поле анкеты | Описание | Банк 1 | Банк 2 | Банк 3 | Обязательность |
|-------------|----------|--------|--------|--------|----------------|
| `income.salary` | Зарплата | `salary` | `base_salary` | `wage` | Обязательно |
| `income.additionalIncome` | Доп. доходы | `additional_income` | `extra_income` | `other_income` | Опционально |
| `income.totalMonthlyIncome` | Общий доход | `total_income` | `monthly_income` | `total_income` | Обязательно |
| `expenses.rent` | Аренда | `rent` | `housing_cost` | `rent` | Опционально |
| `expenses.utilities` | Коммунальные | `utilities` | `utility_bills` | `utilities` | Опционально |
| `expenses.totalMonthlyExpenses` | Общие расходы | `total_expenses` | `monthly_expenses` | `expenses` | Обязательно |
| `property.hasRealEstate` | Недвижимость | `has_property` | `real_estate` | `property` | Обязательно |
| `property.realEstateValue` | Стоимость недвижимости | `property_value` | `real_estate_value` | `property_worth` | Опционально |
| `property.hasVehicle` | Транспорт | `has_vehicle` | `vehicle` | `car` | Обязательно |
| `property.vehicleValue` | Стоимость транспорта | `vehicle_value` | `car_value` | `vehicle_worth` | Опционально |
| `creditHistory.hasActiveLoans` | Активные кредиты | `active_loans` | `current_loans` | `existing_loans` | Обязательно |
| `creditHistory.activeLoansCount` | Количество кредитов | `loans_count` | `number_of_loans` | `loan_count` | Опционально |
| `creditHistory.totalDebt` | Общий долг | `total_debt` | `outstanding_debt` | `debt_amount` | Опционально |
| `creditHistory.hasOverdue` | Просрочки | `has_overdue` | `overdue_payments` | `delinquent` | Обязательно |
| `creditHistory.overdueAmount` | Сумма просрочки | `overdue_amount` | `overdue_sum` | `delinquent_amount` | Опционально |

### 5. Семейные данные (Family Data)

| Поле анкеты | Описание | Банк 1 | Банк 2 | Банк 3 | Обязательность |
|-------------|----------|--------|--------|--------|----------------|
| `spouse.hasSpouse` | Есть супруг/а | `has_spouse` | `married` | `spouse` | Обязательно |
| `spouse.spouseName` | Имя супруга/и | `spouse_name` | `partner_name` | `spouse_name` | Опционально |
| `spouse.spousePhone` | Телефон супруга/и | `spouse_phone` | `partner_phone` | `spouse_phone` | Опционально |
| `spouse.spouseWork` | Работа супруга/и | `spouse_work` | `partner_work` | `spouse_employment` | Опционально |
| `spouse.spouseIncome` | Доход супруга/и | `spouse_income` | `partner_income` | `spouse_income` | Опционально |
| `children[].name` | Имя ребенка | `child_name` | `children[].name` | `child_name` | Опционально |
| `children[].birthDate` | Дата рождения ребенка | `child_birth_date` | `children[].birth_date` | `child_birthday` | Опционально |
| `children[].relationship` | Родство | `child_relationship` | `children[].relation` | `child_relation` | Опционально |
| `emergencyContacts[].name` | Имя контакта | `emergency_name` | `contacts[].name` | `emergency_name` | Обязательно |
| `emergencyContacts[].relationship` | Родство | `emergency_relation` | `contacts[].relation` | `emergency_relation` | Обязательно |
| `emergencyContacts[].phone` | Телефон контакта | `emergency_phone` | `contacts[].phone` | `emergency_phone` | Обязательно |
| `emergencyContacts[].address` | Адрес контакта | `emergency_address` | `contacts[].address` | `emergency_address` | Опционально |

### 6. Дополнительные данные (Additional Data)

| Поле анкеты | Описание | Банк 1 | Банк 2 | Банк 3 | Обязательность |
|-------------|----------|--------|--------|--------|----------------|
| `additionalInfo.hasCriminalRecord` | Судимость | `criminal_record` | `convictions` | `criminal_history` | Обязательно |
| `additionalInfo.hasAdministrativeViolations` | Админ. нарушения | `admin_violations` | `violations` | `admin_offenses` | Обязательно |
| `additionalInfo.hasTaxDebts` | Налоговые долги | `tax_debts` | `tax_liabilities` | `tax_arrears` | Обязательно |
| `additionalInfo.hasAlimonyObligations` | Алименты | `alimony` | `support_obligations` | `child_support` | Обязательно |
| `additionalInfo.additionalComments` | Доп. комментарии | `comments` | `additional_info` | `notes` | Опционально |
| `consents.dataProcessing` | Согласие на обработку | `data_consent` | `privacy_consent` | `data_agreement` | Обязательно |
| `consents.creditHistory` | Согласие на КИ | `credit_consent` | `credit_check` | `credit_agreement` | Обязательно |
| `consents.scoring` | Согласие на скоринг | `scoring_consent` | `scoring_agreement` | `scoring_consent` | Обязательно |
| `consents.marketing` | Маркетинговые рассылки | `marketing_consent` | `marketing_agreement` | `marketing_consent` | Опционально |

## Специфичные поля для разных типов заявок

### Кредит на ПОС (Point of Sale)

| Поле | Описание | Банк 1 | Банк 2 | Банк 3 |
|------|----------|--------|--------|--------|
| `posMerchantId` | ID мерчанта | `merchant_id` | `pos_merchant` | `merchant_code` |
| `posTerminalId` | ID терминала | `terminal_id` | `pos_terminal` | `terminal_code` |
| `posMonthlyTurnover` | Месячный оборот | `monthly_turnover` | `pos_volume` | `terminal_volume` |
| `posAverageTransaction` | Средний чек | `avg_transaction` | `average_ticket` | `avg_ticket` |

### Банковская гарантия

| Поле | Описание | Банк 1 | Банк 2 | Банк 3 |
|------|----------|--------|--------|--------|
| `guaranteeAmount` | Сумма гарантии | `guarantee_amount` | `guarantee_sum` | `guarantee_value` |
| `guaranteePurpose` | Цель гарантии | `guarantee_purpose` | `guarantee_reason` | `guarantee_use` |
| `guaranteePeriod` | Срок гарантии | `guarantee_period` | `guarantee_term` | `guarantee_duration` |
| `guaranteeBeneficiary` | Бенефициар | `beneficiary` | `guarantee_beneficiary` | `guarantee_recipient` |

## Форматы данных

### Даты
- **Стандарт**: ISO 8601 (YYYY-MM-DD)
- **Банк 1**: DD.MM.YYYY
- **Банк 2**: MM/DD/YYYY
- **Банк 3**: YYYY-MM-DD

### Телефоны
- **Стандарт**: +7XXXXXXXXXX
- **Банк 1**: 8XXXXXXXXXX
- **Банк 2**: +7 (XXX) XXX-XX-XX
- **Банк 3**: +7-XXX-XXX-XX-XX

### Суммы
- **Стандарт**: Число с плавающей точкой
- **Банк 1**: Целое число (копейки)
- **Банк 2**: Строка с разделителями
- **Банк 3**: Число с плавающей точкой

## Валидация полей

### Обязательные поля
Все поля, помеченные как "Обязательно", должны быть заполнены для успешной отправки заявки.

### Форматы валидации
- **ИНН**: 10 или 12 цифр
- **СНИЛС**: XXX-XXX-XXX XX
- **Паспорт**: XXXX XXXXXX
- **Телефон**: +7XXXXXXXXXX
- **Email**: RFC 5322

### Ограничения
- **Размер файлов**: Максимум 10MB
- **Количество файлов**: До 10 на заявку
- **Размер текстовых полей**: До 1000 символов
- **Количество детей**: До 10
- **Количество контактов**: До 5

## Примеры трансформации

### Банк 1 (Сбербанк)
```json
{
  "surname": "Иванов",
  "name": "Иван",
  "patronymic": "Иванович",
  "birth_date": "15.03.1985",
  "phone": "89161234567",
  "email": "ivan@example.com",
  "monthly_income": 100000,
  "has_property": true,
  "property_value": 5000000
}
```

### Банк 2 (ВТБ)
```json
{
  "last_name": "Иванов",
  "first_name": "Иван",
  "middle_name": "Иванович",
  "date_of_birth": "03/15/1985",
  "mobile_phone": "+7 (916) 123-45-67",
  "email_address": "ivan@example.com",
  "salary": 100000,
  "real_estate": true,
  "real_estate_value": "5,000,000"
}
```

### Банк 3 (Альфа-Банк)
```json
{
  "family_name": "Иванов",
  "given_name": "Иван",
  "patronymic": "Иванович",
  "birthday": "1985-03-15",
  "contact_phone": "+7-916-123-45-67",
  "email": "ivan@example.com",
  "income": 100000.0,
  "property": true,
  "property_worth": 5000000.0
}
```

## Примечания

1. **Версионирование**: Каждый банк может иметь разные версии API с различными форматами полей.

2. **Локализация**: Некоторые банки требуют поля на английском языке, другие - на русском.

3. **Кодировка**: Все текстовые поля должны быть в кодировке UTF-8.

4. **Обязательность**: Поля, помеченные как "Опционально", могут быть пустыми, но рекомендуется их заполнять для повышения вероятности одобрения.

5. **Валидация**: Перед отправкой в банк все поля проходят дополнительную валидацию согласно требованиям конкретного банка.
