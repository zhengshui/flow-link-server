# FitEasy (健身易) - 后端API接口文档

**版本**: v1.1.0
**更新日期**: 2025-11-03
**基础URL**: `https://api.fiteasy.com` (待定)

---

## 目录

1. [通用说明](#通用说明)
2. [认证接口](#认证接口)
3. [用户接口](#用户接口)
4. [训练记录接口](#训练记录接口)
5. [健身计划接口](#健身计划接口)
6. [计划模板接口](#计划模板接口)
7. [统计数据接口](#统计数据接口)
8. [数据模型](#数据模型)

---

## 通用说明

### 请求头 (Request Headers)

所有需要认证的接口必须包含以下请求头：

```
Content-Type: application/json
Authorization: Bearer {access_token}
```

### 统一响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

**响应码说明**:
- `200`: 成功
- `400`: 请求参数错误
- `401`: 未授权（token无效或过期）
- `403`: 禁止访问
- `404`: 资源不存在
- `500`: 服务器内部错误

---

## 认证接口

### 1. 用户注册

**接口**: `POST /api/auth/register`

**请求参数**:
```json
{
  "username": "string",      // 用户名（4-20字符）
  "password": "string",      // 密码（6-20字符）
  "nickname": "string",      // 昵称（可选）
  "email": "string",         // 邮箱（可选）
  "phone": "string",         // 手机号（可选）
  "gender": "string",        // 性别：男/女（可选）
  "age": 0,                  // 年龄（可选）
  "height": 0,               // 身高cm（可选）
  "weight": 0,               // 体重kg（可选）
  "targetWeight": 0,         // 目标体重kg（可选）
  "fitnessGoal": "string"    // 健身目标：增肌/减脂/力量提升/耐力提升/综合健身（可选）
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "role": "user"
  }
}
```

**说明**: 注册成功后返回访问令牌，客户端应保存token并调用 `GET /api/user/info` 获取完整用户信息

---

### 2. 用户登录

**接口**: `POST /api/auth/login`

**请求参数**:
```json
{
  "username": "string",      // 用户名
  "password": "string"       // 密码
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "role": "user"
  }
}
```

**说明**: 登录成功后返回访问令牌和用户角色（user/admin），客户端应保存token并调用 `GET /api/user/info` 获取完整用户信息

---

### 3. 刷新Token

**接口**: `POST /api/auth/refresh`

**请求头**: `Authorization: Bearer {refresh_token}`

**响应示例**:
```json
{
  "code": 200,
  "message": "刷新成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

## 用户接口

### 1. 获取用户信息

**接口**: `GET /api/user/info`

**需要认证**: 是

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "nickname": "健身达人",
    "avatarUrl": "https://cdn.fiteasy.com/avatar/1.jpg",
    "email": "test@example.com",
    "phone": "13800138000",
    "gender": "男",
    "age": 28,
    "height": 175,
    "weight": 70,
    "targetWeight": 68,
    "fitnessGoal": "增肌",
    "joinDate": "2025-01-01"
  }
}
```

---

### 2. 更新用户信息

**接口**: `PUT /api/user/info`

**需要认证**: 是

**请求参数**:
```json
{
  "nickname": "string",      // 昵称（可选）
  "avatarUrl": "string",     // 头像URL（可选）
  "email": "string",         // 邮箱（可选）
  "phone": "string",         // 手机号（可选）
  "gender": "string",        // 性别：男/女（可选）
  "age": 0,                  // 年龄（可选）
  "height": 0,               // 身高cm（可选）
  "weight": 0,               // 体重kg（可选）
  "targetWeight": 0,         // 目标体重kg（可选）
  "fitnessGoal": "string"    // 健身目标（可选）
}
```

**说明**: 所有字段均为可选，只需传入需要更新的字段

**响应示例**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": null
}
```

---

## 训练记录接口

### 1. 获取训练记录列表

**接口**: `GET /api/training/records`

**需要认证**: 是

**请求参数** (Query):
- `page`: 页码（默认1）
- `pageSize`: 每页条数（默认20）
- `startDate`: 开始日期（可选，格式：YYYY-MM-DD 或 YYYY-MM-DD HH:mm:ss）
- `endDate`: 结束日期（可选，格式：YYYY-MM-DD 或 YYYY-MM-DD HH:mm:ss）
- `planId`: 关联计划ID（可选）

**说明**：
- 如果 `startDate` 使用 YYYY-MM-DD 格式，自动补充为当天 00:00:00
- 如果 `endDate` 使用 YYYY-MM-DD 格式，自动补充为当天 23:59:59

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 50,
    "page": 1,
    "pageSize": 20,
    "records": [
      {
        "id": 1,
        "userId": 1,
        "title": "腿部训练日",
        "startTime": "2025-11-01 09:00:00",
        "endTime": "2025-11-01 10:30:00",
        "duration": 90,
        "exercises": [
          {
            "id": 1,
            "name": "杠铃深蹲",
            "sets": 4,
            "reps": 10,
            "weight": 80,
            "restTime": 90,
            "muscleGroup": "腿部",
            "notes": "保持核心稳定",
            "duration": 20
          }
        ],
        "totalWeight": 5600,
        "totalSets": 11,
        "caloriesBurned": 450,
        "notes": "状态不错",
        "mood": "优秀",
        "planId": 1,
        "createdAt": "2025-11-01 10:30:00",
        "updatedAt": "2025-11-01 10:30:00"
      }
    ]
  }
}
```

---

### 2. 获取单条训练记录

**接口**: `GET /api/training/records/{recordId}`

**需要认证**: 是

**路径参数**:
- `recordId`: 训练记录ID

**响应示例**: 同上单条记录格式

---

### 3. 创建训练记录

**接口**: `POST /api/training/records`

**需要认证**: 是

**请求参数**:
```json
{
  "title": "string",                 // 训练标题（必填）
  "startTime": "string",             // 开始时间 (YYYY-MM-DD HH:mm:ss)（必填）
  "endTime": "string",               // 结束时间 (YYYY-MM-DD HH:mm:ss)（必填）
  "duration": 90,                    // 总时长（分钟）（必填）
  "exercises": [                     // 训练项目列表（必填）
    {
      "name": "string",              // 项目名称
      "sets": 4,                     // 组数
      "reps": 10,                    // 次数
      "weight": 80,                  // 重量（kg）
      "restTime": 90,                // 休息时间（秒）
      "muscleGroup": "string",       // 目标肌群
      "notes": "string",             // 备注
      "duration": 20                 // 训练时长（分钟）
    }
  ],
  "totalWeight": 5600,               // 总重量（kg）（可选）
  "totalSets": 11,                   // 总组数（可选）
  "caloriesBurned": 450,             // 消耗卡路里（可选）
  "notes": "string",                 // 训练备注（可选）
  "mood": "string",                  // 训练状态（优秀/良好/一般/疲劳）（可选）
  "planId": 1                        // 关联计划ID（可选，0表示无计划）
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": 1,
    "createdAt": "2025-11-01 10:30:00"
  }
}
```

---

### 4. 更新训练记录

**接口**: `PUT /api/training/records/{recordId}`

**需要认证**: 是

**路径参数**:
- `recordId`: 训练记录ID

**请求参数**: 同创建训练记录

**响应示例**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": null
}
```

---

### 5. 删除训练记录

**接口**: `DELETE /api/training/records/{recordId}`

**需要认证**: 是

**路径参数**:
- `recordId`: 训练记录ID

**响应示例**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

---

## 健身计划接口

### 1. 获取用户计划列表

**接口**: `GET /api/plans`

**需要认证**: 是

**请求参数** (Query):
- `status`: 计划状态（可选，进行中/已完成/已暂停/已归档）
- `page`: 页码（默认1）
- `pageSize`: 每页条数（默认20）

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 5,
    "page": 1,
    "pageSize": 20,
    "plans": [
      {
        "id": 1,
        "userId": 1,
        "templateId": 1,
        "name": "增肌计划",
        "description": "系统化增肌训练",
        "goal": "增肌",
        "durationWeeks": 8,
        "trainingDaysPerWeek": 4,
        "trainingDays": [
          {
            "dayNumber": 1,
            "dayName": "胸肌+三头日",
            "isRestDay": false,
            "exercises": [...],
            "notes": "重点训练胸大肌中束和上束"
          }
        ],
        "startDate": "2025-10-18",
        "endDate": "2025-12-13",
        "status": "进行中",
        "currentWeek": 3,
        "currentDay": 2,
        "completedDays": [1, 2, 4, 5, 8, 9, 11, 12, 15, 16],
        "totalCompletedDays": 10,
        "completionRate": 27,
        "createdAt": "2025-10-18 10:00:00",
        "updatedAt": "2025-11-01 08:00:00"
      }
    ]
  }
}
```

---

### 2. 获取单个计划详情

**接口**: `GET /api/plans/{planId}`

**需要认证**: 是

**路径参数**:
- `planId`: 计划ID

**响应示例**: 同上单个计划格式

---

### 3. 基于模板创建计划

**接口**: `POST /api/plans/from-template`

**需要认证**: 是

**请求参数**:
```json
{
  "templateId": 1,           // 模板ID
  "startDate": "string",     // 开始日期 (YYYY-MM-DD)
  "name": "string"           // 计划名称（可选，默认使用模板名称）
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": 1,
    "name": "增肌计划",
    "startDate": "2025-11-01",
    "endDate": "2025-12-27"
  }
}
```

---

### 4. 创建自定义计划

**接口**: `POST /api/plans/custom`

**需要认证**: 是

**请求参数**:
```json
{
  "name": "string",                  // 计划名称
  "description": "string",           // 计划描述
  "goal": "string",                  // 训练目标
  "durationWeeks": 8,                // 计划周期（周）
  "trainingDaysPerWeek": 4,          // 每周训练天数
  "trainingDays": [                  // 训练日程
    {
      "dayNumber": 1,
      "dayName": "string",
      "isRestDay": false,
      "exercises": [...],
      "notes": "string"
    }
  ],
  "startDate": "string"              // 开始日期 (YYYY-MM-DD)
}
```

**响应示例**: 同上

---

### 5. 标记训练日为已完成

**接口**: `POST /api/plans/{planId}/complete-day`

**需要认证**: 是

**路径参数**:
- `planId`: 计划ID

**请求参数**:
```json
{
  "dayNumber": 1,            // 第几天
  "recordId": 1              // 关联的训练记录ID（可选）
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "已标记为完成",
  "data": {
    "totalCompletedDays": 11,
    "completionRate": 30
  }
}
```

---

### 6. 更新计划状态

**接口**: `PUT /api/plans/{planId}/status`

**需要认证**: 是

**路径参数**:
- `planId`: 计划ID

**请求参数**:
```json
{
  "status": "string"         // 进行中/已完成/已暂停/已归档
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": null
}
```

---

### 7. 删除计划

**接口**: `DELETE /api/plans/{planId}`

**需要认证**: 是

**路径参数**:
- `planId`: 计划ID

**响应示例**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

---

## 计划模板接口

### 1. 获取模板列表

**接口**: `GET /api/templates`

**需要认证**: 否（公开接口）

**请求参数** (Query):
- `goal`: 训练目标（可选，增肌/减脂/力量提升/耐力提升/综合健身）
- `level`: 难度等级（可选，初级/中级/高级）
- `page`: 页码（默认1）
- `pageSize`: 每页条数（默认20）

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 10,
    "page": 1,
    "pageSize": 20,
    "templates": [
      {
        "id": 1,
        "name": "增肌计划",
        "description": "系统化增肌训练，适合有一定基础的健身者",
        "goal": "增肌",
        "level": "中级",
        "durationWeeks": 8,
        "trainingDaysPerWeek": 4,
        "trainingDays": [...],
        "imageUrl": "https://cdn.fiteasy.com/templates/1.jpg",
        "author": "FitEasy官方",
        "tags": ["增肌", "中级", "器械训练"],
        "createdAt": "2025-01-01 10:00:00"
      }
    ]
  }
}
```

---

### 2. 获取单个模板详情

**接口**: `GET /api/templates/{templateId}`

**需要认证**: 否

**路径参数**:
- `templateId`: 模板ID

**响应示例**: 同上单个模板格式

---

## 统计数据接口

### 1. 获取训练统计数据

**接口**: `GET /api/stats/training`

**需要认证**: 是

**请求参数** (Query):
- `period`: 统计周期（week/month/year）
- `startDate`: 开始日期（可选，格式：YYYY-MM-DD）
- `endDate`: 结束日期（可选，格式：YYYY-MM-DD）

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "userId": 1,
    "period": "week",
    "startDate": "2025-10-26",
    "endDate": "2025-11-01",
    "totalTrainingCount": 5,
    "totalDuration": 350,
    "totalWeight": 16400,
    "totalSets": 55,
    "totalCalories": 1850,
    "avgDuration": 70,
    "avgWeight": 3280,
    "mostTrainedMuscle": "胸部",
    "favoriteExercise": "杠铃卧推",
    "dailyStats": [
      {
        "date": "2025-10-26",
        "trainingCount": 0,
        "duration": 0,
        "weight": 0,
        "sets": 0,
        "calories": 0
      },
      {
        "date": "2025-10-27",
        "trainingCount": 1,
        "duration": 75,
        "weight": 4600,
        "sets": 11,
        "calories": 380
      }
    ]
  }
}
```

---

### 2. 获取肌群训练统计

**接口**: `GET /api/stats/muscle-groups`

**需要认证**: 是

**请求参数** (Query):
- `period`: 统计周期（week/month/year，默认month）

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "period": "month",
    "muscleGroups": [
      {
        "muscleGroup": "胸部",
        "trainingCount": 12,
        "totalWeight": 15000,
        "percentage": 25
      },
      {
        "muscleGroup": "背部",
        "trainingCount": 10,
        "totalWeight": 12000,
        "percentage": 20
      },
      {
        "muscleGroup": "腿部",
        "trainingCount": 14,
        "totalWeight": 18000,
        "percentage": 30
      }
    ]
  }
}
```

---

### 3. 获取个人记录 (PR - Personal Record)

**接口**: `GET /api/stats/personal-records`

**需要认证**: 是

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "records": [
      {
        "exerciseName": "杠铃深蹲",
        "maxWeight": 120,
        "date": "2025-10-15",
        "recordId": 15
      },
      {
        "exerciseName": "杠铃卧推",
        "maxWeight": 80,
        "date": "2025-10-20",
        "recordId": 18
      },
      {
        "exerciseName": "硬拉",
        "maxWeight": 140,
        "date": "2025-10-25",
        "recordId": 22
      }
    ]
  }
}
```

---

### 4. 获取训练日历

**接口**: `GET /api/stats/calendar`

**需要认证**: 是

**请求参数** (Query):
- `year`: 年份（默认当前年）
- `month`: 月份（默认当前月，1-12）

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "year": 2025,
    "month": 11,
    "days": [
      {
        "date": "2025-11-01",
        "hasTraining": true,
        "trainingCount": 1,
        "totalDuration": 90
      },
      {
        "date": "2025-11-02",
        "hasTraining": false,
        "trainingCount": 0,
        "totalDuration": 0
      }
    ]
  }
}
```

---

## 数据模型

### TrainingRecord (训练记录)

```typescript
{
  id: number                      // 训练记录ID
  userId: number                  // 用户ID
  title: string                   // 训练标题
  date: string                    // 训练日期 (YYYY-MM-DD) - 从 startTime 提取
  startTime: string               // 开始时间 (YYYY-MM-DD HH:mm:ss) 完整日期时间
  endTime: string                 // 结束时间 (YYYY-MM-DD HH:mm:ss) 完整日期时间
  duration: number                // 总时长（分钟）- 由后端根据 startTime 和 endTime 计算
  exercises: Exercise[]           // 训练项目列表
  totalWeight: number             // 总重量（kg）
  totalSets: number               // 总组数
  caloriesBurned: number          // 消耗卡路里
  notes: string                   // 训练备注
  mood: string                    // 训练状态（优秀/良好/一般/疲劳）
  planId: number                  // 关联计划ID（0表示无计划）
  createdAt: string               // 创建时间
  updatedAt: string               // 更新时间
}
```

### Exercise (训练项目)

```typescript
{
  id: number                      // 训练项目ID
  name: string                    // 项目名称
  sets: number                    // 组数
  reps: number                    // 次数
  weight: number                  // 重量（kg）
  restTime: number                // 休息时间（秒）
  muscleGroup: string             // 目标肌群
  notes: string                   // 备注
  duration: number                // 训练时长（分钟）
}
```

### FitnessPlan (健身计划)

```typescript
{
  id: number                      // 计划ID
  userId: number                  // 用户ID
  templateId: number              // 模板ID（0表示自定义）
  name: string                    // 计划名称
  description: string             // 计划描述
  goal: string                    // 训练目标
  durationWeeks: number           // 计划周期（周）
  trainingDaysPerWeek: number     // 每周训练天数
  trainingDays: TrainingDay[]     // 训练日程
  startDate: string               // 开始日期 (YYYY-MM-DD)
  endDate: string                 // 结束日期 (YYYY-MM-DD)
  status: string                  // 计划状态（进行中/已完成/已暂停/已归档）
  currentWeek: number             // 当前第几周
  currentDay: number              // 当前第几天
  completedDays: number[]         // 已完成的训练日
  totalCompletedDays: number      // 累计完成天数
  completionRate: number          // 完成率（百分比）
  createdAt: string               // 创建时间
  updatedAt: string               // 更新时间
}
```

### PlanTemplate (计划模板)

```typescript
{
  id: number                      // 模板ID
  name: string                    // 模板名称
  description: string             // 模板描述
  goal: string                    // 训练目标
  level: string                   // 难度等级（初级/中级/高级）
  durationWeeks: number           // 计划周期（周）
  trainingDaysPerWeek: number     // 每周训练天数
  trainingDays: TrainingDay[]     // 训练日程
  imageUrl: string                // 封面图片URL
  author: string                  // 作者/来源
  tags: string[]                  // 标签
  createdAt: string               // 创建时间
}
```

### TrainingDay (训练日程)

```typescript
{
  dayNumber: number               // 第几天
  dayName: string                 // 训练日名称
  isRestDay: boolean              // 是否为休息日
  exercises: Exercise[]           // 当日训练项目
  notes: string                   // 当日备注
}
```

### TrainingStats (训练统计)

```typescript
{
  userId: number                  // 用户ID
  period: string                  // 统计周期（week/month/year）
  startDate: string               // 统计开始日期
  endDate: string                 // 统计结束日期
  totalTrainingCount: number      // 总训练次数
  totalDuration: number           // 总训练时长（分钟）
  totalWeight: number             // 总重量（kg）
  totalSets: number               // 总组数
  totalCalories: number           // 总消耗卡路里
  avgDuration: number             // 平均训练时长
  avgWeight: number               // 平均单次重量
  mostTrainedMuscle: string       // 训练最多的肌群
  favoriteExercise: string        // 最常做的训练项目
  dailyStats: DailyStats[]        // 每日统计数据
}
```

### DailyStats (每日统计)

```typescript
{
  date: string                    // 日期 (YYYY-MM-DD)
  trainingCount: number           // 当日训练次数
  duration: number                // 当日训练时长
  weight: number                  // 当日总重量
  sets: number                    // 当日总组数
  calories: number                // 当日消耗卡路里
}
```

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权（需要登录） |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 409 | 资源冲突（如用户名已存在） |
| 500 | 服务器内部错误 |

---

## 附录

### 训练状态枚举

- `优秀` - 状态极佳，完成所有训练目标
- `良好` - 状态良好，大部分完成
- `一般` - 普通状态，基本完成
- `疲劳` - 感觉疲劳，未完全完成

### 肌群类别

- `胸部` - 胸大肌、胸小肌
- `背部` - 背阔肌、斜方肌、竖脊肌
- `腿部` - 股四头肌、腘绳肌、臀大肌
- `肩部` - 三角肌
- `手臂` - 肱二头肌、肱三头肌
- `核心` - 腹直肌、腹斜肌
- `有氧` - 跑步、游泳等有氧运动

### 训练目标

- `增肌` - 增加肌肉质量
- `减脂` - 减少体脂率
- `力量提升` - 提高最大力量
- `耐力提升` - 提高心肺耐力
- `综合健身` - 全面身体素质提升

---

## 版本历史

### v1.1.0 (2025-11-03)

**重要更新** - 用户模型完善

1. **用户注册接口增强** (`POST /api/auth/register`)
   - 新增可选字段：phone, gender, age, height, weight, targetWeight, fitnessGoal
   - 支持注册时填写完整的个人信息和健身目标

2. **认证响应格式优化**
   - 登录和注册接口统一返回 `{token, role}` 格式
   - 移除冗余的用户信息字段，使用 `GET /api/user/info` 获取完整信息
   - 提高安全性和接口设计一致性

3. **用户信息模型完善** (`GET /api/user/info`)
   - 新增字段：phone（手机号）、gender（性别）、age（年龄）
   - 新增字段：height（身高cm）、weight（体重kg）、targetWeight（目标体重kg）
   - 新增字段：fitnessGoal（健身目标）、joinDate（加入日期）
   - 将 createdAt 改为 joinDate，语义更明确

4. **用户信息更新接口增强** (`PUT /api/user/info`)
   - 支持更新所有个人信息字段（除id和username外）
   - 包括：身体数据、健身目标等关键健身应用字段

**向后兼容性**:
- 注册接口新增字段均为可选，不影响现有客户端
- 用户信息更新接口新增字段均为可选
- 建议客户端尽快升级以支持完整的用户信息管理

---

### v1.0.0 (2025-11-01)

- 初始版本，包含核心功能接口

---

**联系方式**: api-support@fiteasy.com
