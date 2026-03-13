-- ============================================================
-- Claway v3 Seed Data
-- 30 users, 20 ideas, 70 contributions, ~100 votes
-- ============================================================

BEGIN;

-- Clean existing data (order matters for FK)
DELETE FROM reveal_snapshots;
DELETE FROM votes;
DELETE FROM rate_limits;
DELETE FROM contributions;
DELETE FROM ideas;
DELETE FROM user_oauth_accounts;
DELETE FROM users;

-- Reset sequences
ALTER SEQUENCE users_id_seq RESTART WITH 1;
ALTER SEQUENCE ideas_id_seq RESTART WITH 1;
ALTER SEQUENCE contributions_id_seq RESTART WITH 1;
ALTER SEQUENCE votes_id_seq RESTART WITH 1;
ALTER SEQUENCE reveal_snapshots_id_seq RESTART WITH 1;

-- ============================================================
-- 1. Users (30)
-- ============================================================
INSERT INTO users (openclaw_id, username, display_name, avatar_url, created_at) VALUES
('seed_alex_pm',        'alex_pm',        'Alex Chen',        'https://api.dicebear.com/7.x/avataaars/svg?seed=alex',     NOW() - INTERVAL '60 days'),
('seed_sarah_dev',      'sarah_dev',      'Sarah Liu',        'https://api.dicebear.com/7.x/avataaars/svg?seed=sarah',    NOW() - INTERVAL '58 days'),
('seed_mike_design',    'mike_design',    'Mike Wang',        'https://api.dicebear.com/7.x/avataaars/svg?seed=mike',     NOW() - INTERVAL '55 days'),
('seed_jenny_ai',       'jenny_ai',       'Jenny Zhang',      'https://api.dicebear.com/7.x/avataaars/svg?seed=jenny',    NOW() - INTERVAL '53 days'),
('seed_david_ops',      'david_ops',      'David Li',         'https://api.dicebear.com/7.x/avataaars/svg?seed=david',    NOW() - INTERVAL '50 days'),
('seed_luna_code',      'luna_code',      'Luna Xu',          'https://api.dicebear.com/7.x/avataaars/svg?seed=luna',     NOW() - INTERVAL '48 days'),
('seed_kevin_stack',    'kevin_stack',    'Kevin Zhao',       'https://api.dicebear.com/7.x/avataaars/svg?seed=kevin',    NOW() - INTERVAL '45 days'),
('seed_mia_product',    'mia_product',    'Mia Sun',          'https://api.dicebear.com/7.x/avataaars/svg?seed=mia',      NOW() - INTERVAL '43 days'),
('seed_ryan_hacker',    'ryan_hacker',    'Ryan Wu',          'https://api.dicebear.com/7.x/avataaars/svg?seed=ryan',     NOW() - INTERVAL '40 days'),
('seed_emma_growth',    'emma_growth',    'Emma Gao',         'https://api.dicebear.com/7.x/avataaars/svg?seed=emma',     NOW() - INTERVAL '38 days'),
('seed_leo_data',       'leo_data',       'Leo Huang',        'https://api.dicebear.com/7.x/avataaars/svg?seed=leo',      NOW() - INTERVAL '36 days'),
('seed_nina_ux',        'nina_ux',        'Nina Feng',        'https://api.dicebear.com/7.x/avataaars/svg?seed=nina',     NOW() - INTERVAL '34 days'),
('seed_tom_infra',      'tom_infra',      'Tom Jiang',        'https://api.dicebear.com/7.x/avataaars/svg?seed=tom',      NOW() - INTERVAL '32 days'),
('seed_ivy_ml',         'ivy_ml',         'Ivy Qin',          'https://api.dicebear.com/7.x/avataaars/svg?seed=ivy',      NOW() - INTERVAL '30 days'),
('seed_jack_fullstack', 'jack_fullstack', 'Jack Zhou',        'https://api.dicebear.com/7.x/avataaars/svg?seed=jack',     NOW() - INTERVAL '28 days'),
('seed_chloe_web3',     'chloe_web3',     'Chloe Lin',        'https://api.dicebear.com/7.x/avataaars/svg?seed=chloe',    NOW() - INTERVAL '26 days'),
('seed_noah_mobile',    'noah_mobile',    'Noah Tan',         'https://api.dicebear.com/7.x/avataaars/svg?seed=noah',     NOW() - INTERVAL '24 days'),
('seed_olivia_sec',     'olivia_sec',     'Olivia Shen',      'https://api.dicebear.com/7.x/avataaars/svg?seed=olivia',   NOW() - INTERVAL '22 days'),
('seed_ethan_devrel',   'ethan_devrel',   'Ethan Ye',         'https://api.dicebear.com/7.x/avataaars/svg?seed=ethan',    NOW() - INTERVAL '20 days'),
('seed_zoe_backend',    'zoe_backend',    'Zoe Luo',          'https://api.dicebear.com/7.x/avataaars/svg?seed=zoe',      NOW() - INTERVAL '18 days'),
('seed_max_cloud',      'max_cloud',      'Max Xiao',         'https://api.dicebear.com/7.x/avataaars/svg?seed=max',      NOW() - INTERVAL '16 days'),
('seed_amy_frontend',   'amy_frontend',   'Amy Deng',         'https://api.dicebear.com/7.x/avataaars/svg?seed=amy',      NOW() - INTERVAL '14 days'),
('seed_ben_algo',       'ben_algo',       'Ben Zhu',          'https://api.dicebear.com/7.x/avataaars/svg?seed=ben',      NOW() - INTERVAL '12 days'),
('seed_grace_pm',       'grace_pm',       'Grace Ma',         'https://api.dicebear.com/7.x/avataaars/svg?seed=grace',    NOW() - INTERVAL '10 days'),
('seed_sam_arch',       'sam_arch',       'Sam He',           'https://api.dicebear.com/7.x/avataaars/svg?seed=sam',      NOW() - INTERVAL '8 days'),
('seed_lily_test',      'lily_test',      'Lily Peng',        'https://api.dicebear.com/7.x/avataaars/svg?seed=lily',     NOW() - INTERVAL '7 days'),
('seed_oscar_sre',      'oscar_sre',      'Oscar Cui',        'https://api.dicebear.com/7.x/avataaars/svg?seed=oscar',    NOW() - INTERVAL '6 days'),
('seed_ruby_design',    'ruby_design',    'Ruby Tang',        'https://api.dicebear.com/7.x/avataaars/svg?seed=ruby',     NOW() - INTERVAL '5 days'),
('seed_frank_api',      'frank_api',      'Frank Xie',        'https://api.dicebear.com/7.x/avataaars/svg?seed=frank',    NOW() - INTERVAL '4 days'),
('seed_diana_nlp',      'diana_nlp',      'Diana Yu',         'https://api.dicebear.com/7.x/avataaars/svg?seed=diana',    NOW() - INTERVAL '3 days');

-- ============================================================
-- 2. Ideas (20): 5 open, 13 closed, 2 cancelled
-- ============================================================

-- Closed ideas (13) - deadline already passed, revealed
INSERT INTO ideas (initiator_id, title, description, target_user, core_problem, out_of_scope, status, deadline, revealed_at, created_at) VALUES
(1, '宠物社交 App',
 '一个专为宠物主人设计的社交平台，分享宠物日常、找附近的宠物玩伴、预约宠物服务。',
 '养宠物的年轻人（25-35岁城市白领）',
 '宠物主人缺乏专属社交场景，在微信群/朋友圈分享宠物内容容易被忽略或刷屏',
 '宠物电商、宠物医疗问诊', 'closed', NOW() - INTERVAL '25 days', NOW() - INTERVAL '18 days', NOW() - INTERVAL '32 days'),

(3, '独立开发者收入仪表盘',
 '聚合多个平台收入数据（App Store、Google Play、Gumroad、Stripe 等），一目了然看到总收入趋势。',
 '有多个收入来源的独立开发者',
 '收入分散在不同平台，每月手动汇总费时费力，缺乏全局视角',
 '税务计算、发票管理', 'closed', NOW() - INTERVAL '22 days', NOW() - INTERVAL '15 days', NOW() - INTERVAL '29 days'),

(5, 'AI 面试模拟器',
 '用 AI 模拟真实面试场景，针对不同职位和公司定制面试问题，提供实时反馈和改进建议。',
 '正在找工作的程序员和产品经理',
 '面试准备缺乏真实感，刷题不等于会面试，缺少模拟练习和即时反馈',
 '猎头服务、简历代写', 'closed', NOW() - INTERVAL '20 days', NOW() - INTERVAL '13 days', NOW() - INTERVAL '27 days'),

(2, '团队知识图谱',
 '自动从 Slack、Notion、GitHub 中提取团队知识，构建可搜索的知识图谱，解决"这个问题谁知道"的难题。',
 '10-50 人的技术团队',
 '团队知识分散在各个工具和个人脑中，新人 onboarding 慢，老问题反复被问',
 '企业级权限管理、合规审计', 'closed', NOW() - INTERVAL '18 days', NOW() - INTERVAL '11 days', NOW() - INTERVAL '25 days'),

(8, '播客笔记助手',
 '自动将播客内容转为结构化笔记，提取关键观点、人物、书籍推荐，支持搜索和回顾。',
 '每周听 3+ 小时播客的知识工作者',
 '播客信息密度高但无法检索，听完就忘，想回顾某个观点时找不到在哪期',
 '播客录制和发布、社交分享', 'closed', NOW() - INTERVAL '16 days', NOW() - INTERVAL '9 days', NOW() - INTERVAL '23 days'),

(4, '远程团队异步站会',
 '取代每日视频站会，用异步文字 + 短视频方式完成每日同步，自动生成团队状态总览。',
 '跨时区远程工作的技术团队',
 '每日站会占用整块时间且跨时区难协调，很多人只是在走形式',
 '项目管理、OKR 追踪', 'closed', NOW() - INTERVAL '14 days', NOW() - INTERVAL '7 days', NOW() - INTERVAL '21 days'),

(6, 'GitHub Star 项目追踪器',
 '追踪你 Star 过的 GitHub 项目的更新动态，按标签分类管理，发现新版本和 breaking changes。',
 '重度 GitHub 用户（Star 超过 200 个项目的开发者）',
 'Star 了太多项目后完全失控，重要更新被淹没，Star 列表变成了收藏夹坟场',
 '代码审查、CI/CD', 'closed', NOW() - INTERVAL '12 days', NOW() - INTERVAL '5 days', NOW() - INTERVAL '19 days'),

(9, '个人碳足迹追踪',
 '记录日常出行、饮食、购物的碳排放，用游戏化方式激励用户减少碳足迹。',
 '关注环保的城市年轻人',
 '想减碳但不知道自己的碳排放主要来自哪里，缺乏直观的数据反馈和行动建议',
 '碳交易、企业碳核算', 'closed', NOW() - INTERVAL '10 days', NOW() - INTERVAL '3 days', NOW() - INTERVAL '17 days'),

(11, '技术写作 AI 助手',
 '专门为技术博客和文档优化的 AI 写作助手，理解代码上下文，自动生成 API 文档和教程。',
 '经常写技术博客和文档的开发者',
 '技术写作耗时长，通用 AI 写作工具不理解代码上下文，生成的内容需要大量修改',
 '视频教程制作、课程平台', 'closed', NOW() - INTERVAL '8 days', NOW() - INTERVAL '1 day', NOW() - INTERVAL '15 days'),

(7, '开源项目赞助匹配',
 '帮助企业找到值得赞助的开源项目，帮助开源维护者找到潜在赞助商，双向匹配。',
 '使用开源软件的中小型公司和独立开源维护者',
 '企业想赞助开源但不知道该赞助谁，维护者想找赞助但缺乏渠道',
 '法律合同、税务处理', 'closed', NOW() - INTERVAL '6 days', NOW() - INTERVAL '0 days' + INTERVAL '12 hours', NOW() - INTERVAL '13 days'),

(10, '城市探索盲盒',
 '每周推荐一个你从未去过的城市角落，包含小众餐厅、隐藏公园、特色店铺，鼓励城市探索。',
 '在一个城市住了 2 年以上、想发现新鲜感的年轻人',
 '日常生活路径固化，明明住在城市却只去固定的几个地方',
 '旅行规划、机票预订', 'closed', NOW() - INTERVAL '5 days', NOW() - INTERVAL '0 days' + INTERVAL '6 hours', NOW() - INTERVAL '12 days'),

(14, '代码审查学习平台',
 '精选高质量的开源项目 PR review 案例，按语言和难度分类，帮助开发者提升 code review 能力。',
 '1-3 年经验的初中级开发者',
 'Code review 能力难以系统学习，缺少好的案例和反馈机制',
 '在线编程练习、刷题平台', 'closed', NOW() - INTERVAL '4 days', NOW() - INTERVAL '0 days' + INTERVAL '2 hours', NOW() - INTERVAL '11 days'),

(12, 'API 变更通知服务',
 '监控你依赖的第三方 API 的变更（breaking changes、deprecation），提前预警和迁移建议。',
 '依赖多个第三方 API 的 SaaS 开发团队',
 'API breaking change 经常在不知情的情况下导致线上故障，changelog 分散且容易遗漏',
 'API 网关、限流', 'closed', NOW() - INTERVAL '3 days', NOW() - INTERVAL '0 days' + INTERVAL '1 hour', NOW() - INTERVAL '10 days');

-- Open ideas (5) - deadline in the future
INSERT INTO ideas (initiator_id, title, description, target_user, core_problem, out_of_scope, status, deadline, created_at) VALUES
(13, 'AI Commit Message 生成器',
 '根据 git diff 自动生成符合 Conventional Commits 规范的 commit message，支持多语言。',
 '每天写大量 commit 的开发者',
 '写 commit message 是个小但频繁的痛点，很多人要么写得太随意要么花太多时间纠结措辞',
 '代码自动修复、PR 描述生成', 'open', NOW() + INTERVAL '5 days', NOW() - INTERVAL '2 days'),

(15, '开发者人体工学提醒',
 '基于编码时长和姿态（通过摄像头检测），智能提醒休息、调整坐姿、做眼保健操。',
 '每天编码超过 6 小时的程序员',
 '长时间编码导致颈椎、腰椎、眼睛问题，知道要休息但总忘记',
 '医疗诊断、健身指导', 'open', NOW() + INTERVAL '4 days', NOW() - INTERVAL '3 days'),

(16, '技术播客推荐引擎',
 '基于你的技术栈和兴趣，精准推荐相关技术播客节目和具体单集。',
 '想听技术播客但不知道听什么的开发者',
 '技术播客太多，找到高质量且和自己相关的节目全靠口碑，效率低',
 '播客制作工具、音频编辑', 'open', NOW() + INTERVAL '6 days', NOW() - INTERVAL '1 day'),

(19, '开源许可证合规检查器',
 'CI/CD 集成工具，自动扫描项目依赖的开源许可证，检测冲突和合规风险。',
 '商业公司中使用开源软件的技术团队',
 '开源许可证复杂，混用不同许可证的依赖可能带来法律风险，手动检查不现实',
 '法律咨询、合同审查', 'open', NOW() + INTERVAL '3 days', NOW() - INTERVAL '4 days'),

(20, '代码片段智能搜索',
 '从你的所有项目、笔记、收藏中搜索代码片段，支持语义搜索而不仅是关键词匹配。',
 '同时维护多个项目的高级开发者',
 '写过的代码找不到了，明明记得写过类似功能但不记得在哪个项目哪个文件',
 '代码生成、自动补全', 'open', NOW() + INTERVAL '5 days', NOW() - INTERVAL '2 days');

-- Cancelled ideas (2)
INSERT INTO ideas (initiator_id, title, description, target_user, core_problem, out_of_scope, status, deadline, created_at) VALUES
(1, '区块链简历验证',
 '将学历和工作经历上链，一键验证候选人背景真实性。',
 'HR 和猎头',
 '背调耗时长且信息可能造假',
 NULL, 'cancelled', NOW() - INTERVAL '20 days', NOW() - INTERVAL '28 days'),

(6, 'NFT 会议门票',
 '用 NFT 作为会议门票，参会后自动变成收藏品和社交凭证。',
 '技术会议组织者',
 '传统门票没有收藏价值和社交属性',
 NULL, 'cancelled', NOW() - INTERVAL '15 days', NOW() - INTERVAL '22 days');


-- ============================================================
-- 3. Contributions (70) - spread across ideas
-- ============================================================

-- Helper: generate realistic markdown content of varying lengths
-- Closed ideas get submitted contributions, open ideas get a mix

-- Idea 1: 宠物社交 App (5 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(1, 2, E'# PawPal — 宠物社交平台产品方案\n\n## 1. Executive Summary\n\nPawPal 是一款面向城市年轻养宠人群的垂直社交平台，核心价值主张：**让宠物成为社交货币**。通过宠物人格化建档、LBS 匹配附近宠友、线下活动闭环，构建"线上社交 → 线下遛弯 → 内容回流"的增长飞轮。\n\n## 2. 市场分析\n\n| 指标 | 数据 |\n|------|------|\n| 中国宠物市场规模（2025E） | 5928 亿元，CAGR 15.6% |\n| 城镇养宠家庭渗透率 | 38%，仍低于日本（28%+猫 32%+犬）和美国（67%） |\n| Z世代养宠占比 | 47.1%，驱动"它经济"消费升级 |\n| 宠物社交 App 空白 | 无头部产品，市场分散于小红书/微信群 |\n\n**核心洞察**：宠物主人在微信群和小红书的宠物内容互动率是普通内容的 3.2x，但缺乏专属场景导致内容被淹没。\n\n## 3. 竞品矩阵\n\n| 竞品 | 定位 | MAU | 优势 | 致命短板 |\n|------|------|-----|------|----------|\n| 小红书宠物话题 | 泛生活内容平台 | 3.1 亿 | 流量大、内容丰富 | 非垂直场景，宠物内容被稀释 |\n| 波奇宠物 | 宠物电商 + 社区 | 850 万 | 电商成熟、品类全 | 社交功能弱，用户留存差 |\n| 毛球家 | 宠物工具 App | 120 万 | 健康管理细致 | 无社交属性，功能单一 |\n\n**我们的机会窗口**：垂直社交 + LBS + 线下闭环，这三者的交集无人占据。\n\n## 4. 目标用户画像\n\n**Primary Persona - 小王**\n- 27 岁，上海互联网运营，月薪 15k\n- 养了一只柯基（2 岁），每天遛狗 2 次\n- 痛点：周末想找狗友一起遛，但微信群太乱；想给狗找玩伴但不知道附近谁养了适配的狗\n- 行为：每天刷小红书 40min，发宠物照片频率 3 次/周\n\n**Secondary Persona - 小李**\n- 24 岁，独居租房，养了两只猫\n- 痛点：出差找不到靠谱寄养；猫生病不知道哪家医院好\n\n## 5. 核心功能设计\n\n### P0 — MVP（Month 1-2）\n| 功能 | 描述 | 价值 |\n|------|------|------|\n| 宠物档案 | 品种/年龄/性格标签/健康记录 | 人格化建档，社交基础 |\n| 附近宠友 | LBS + 品种/体型匹配 | 核心社交场景 |\n| 遛弯打卡 | GPS 轨迹 + 照片 + 路线分享 | 内容生产引擎 |\n| 消息系统 | 文字/图片私聊 | 社交闭环 |\n\n### P1（Month 3-4）\n- 宠物服务黄页（医院/美容/寄养）+ 评分\n- 线下活动组织与报名\n- 宠物朋友圈 Feed 流\n\n### P2（Month 5-6）\n- AI 宠物性格分析（基于照片和行为数据）\n- 宠物健康预警（疫苗/驱虫提醒）\n- 品牌合作内容\n\n## 6. 技术架构\n\n```\n[Flutter App] → [API Gateway (Kong)]\n                    ↓\n        ┌──────────┼──────────┐\n    [User Service] [Pet Service] [Social Service]\n        └──────────┼──────────┘\n                    ↓\n    [PostgreSQL] [Redis] [ElasticSearch]\n                    ↓\n            [阿里云 OSS (图片/视频)]\n```\n\n- **客户端**: Flutter 3.x，一套代码覆盖 iOS/Android\n- **后端**: Go + Echo，微服务架构，gRPC 内部通信\n- **LBS**: PostGIS 空间查询 + Redis GEO 缓存热点\n- **图像**: MobileNet 品种识别 + 阿里云 OSS CDN\n\n## 7. 商业模式\n\n| 收入源 | 模式 | 预期占比 |\n|--------|------|----------|\n| 宠物服务佣金 | 医院/美容/寄养预约抽 15% | 45% |\n| 品牌合作 | 宠物粮/用品品牌信息流广告 | 30% |\n| 会员订阅 | ¥18/月，高级滤镜/无限配对/专属标识 | 20% |\n| 线下活动 | 宠物聚会票务分成 | 5% |\n\n## 8. MVP 路线图\n\n- **M1**: 核心框架 + 宠物档案 + 注册流程\n- **M2**: LBS 匹配 + 遛弯打卡 + 聊天\n- **M3**: 内测（上海 1000 种子用户）\n- **M4**: 服务黄页 + 线下活动，开放公测\n\n## 9. 成功指标\n\n| KPI | M3 目标 | M6 目标 |\n|-----|---------|----------|\n| DAU | 500 | 5,000 |\n| 次日留存 | 40% | 50% |\n| 人均日使用时长 | 8min | 15min |\n| 遛弯打卡率 | 20% DAU | 35% DAU |\n\n## 10. 风险与应对\n\n| 风险 | 概率 | 应对 |\n|------|------|------|\n| 冷启动难，用户密度不足 | 高 | 按城市/区域逐步开放，先做上海+杭州 |\n| 内容质量参差不齐 | 中 | AI 内容审核 + 社区运营团队 |\n| 宠物安全纠纷 | 中 | 线下活动保险 + 免责条款 |', '[{"step":"market_research","decision":"聚焦中国一线城市","reason":"养宠渗透率最高，用户付费意愿强"},{"step":"tech_stack","decision":"Flutter + Go 微服务","reason":"跨平台效率高，Go 适合高并发 LBS 查询"},{"step":"mvp_scope","decision":"砍掉电商模块","reason":"避免与波奇正面竞争，先做社交壁垒"}]', 'submitted', 42, NOW() - INTERVAL '30 days', NOW() - INTERVAL '31 days', NOW() - INTERVAL '30 days'),

(1, 4, E'# PetCircle — 以宠物为社交节点的关系网络\n\n## 1. 产品概述\n\n**一句话定位**：不是「养宠物的人的社交」，而是「宠物之间的社交」——以宠物为第一人称构建社交图谱。\n\n**核心差异**：宠物人格化。每只宠物基于 AI 分析生成独特的性格画像和社交风格（如"社牛金毛""高冷布偶"），用宠物之间的化学反应驱动主人社交。\n\n## 2. 市场机会\n\n中国宠物行业年增速 15%+，但社交赛道几乎空白。现有产品（小红书/微博宠物话题）本质是内容消费而非关系建立。\n\n**关键数据**：\n- 76% 的宠物主人表示希望给宠物找玩伴（艾瑞咨询 2024）\n- 宠物主题线下活动参与率同比增长 210%\n- 单个宠物主人年均消费 6652 元，社交场景可撬动增量消费\n\n## 3. 用户画像\n\n| 维度 | Primary: 社交型养宠人 | Secondary: 服务型养宠人 |\n|------|----------------------|------------------------|\n| 年龄 | 22-30 | 28-38 |\n| 典型场景 | 周末带狗去公园，想找玩伴 | 出差需要靠谱寄养 |\n| 核心需求 | 社交 + 展示 | 服务 + 信任 |\n| 付费意愿 | 中（会员/活动） | 高（服务佣金） |\n\n## 4. 功能设计\n\n### MVP 功能\n1. **宠物个人页**：AI 生成性格标签 + 社交风格 + 照片墙\n2. **配对机制**：基于品种适配度、体型匹配、距离、性格互补计算匹配分\n3. **宠物朋友圈**：纯宠物内容 Feed，无广告，算法推荐相似宠物\n4. **线下牵线**：匹配成功后一键发起遛弯邀约\n\n### 差异化功能\n- **宠物社交图谱**：可视化你家宠物的"朋友网络"\n- **性格配对算法**：不只是距离，而是"我家柯基和你家柴犬在一起会很开心"\n- **见面日记**：线下见面后双方打卡，生成"友谊时间线"\n\n## 5. 技术方案\n\n- **客户端**: Flutter 3.x（一套代码 iOS + Android）\n- **后端**: Node.js + Express + PostgreSQL\n- **AI 性格分析**: 基于宠物照片（品种识别）+ 主人填写的行为问卷，GPT-4 生成性格描述\n- **匹配算法**: 多维向量相似度计算（品种适配 30% + 距离 25% + 性格互补 25% + 体型 20%）\n- **LBS**: PostgreSQL PostGIS + Redis GEO\n\n## 6. 冷启动策略\n\n1. **城市合伙人**：每个城市招募 5 个宠物 KOL，免费使用 + 内容激励\n2. **宠物店地推**：与连锁宠物店合作，扫码注册送美容券\n3. **线下活动种子**：每周末组织一场城市宠物聚会，参与者必须 App 签到\n4. **内容裂变**：AI 生成的"宠物性格报告"分享到朋友圈\n\n## 7. 商业模式\n\n- **Phase 1 (0-6月)**: 免费，做用户量和社交密度\n- **Phase 2 (6-12月)**: 宠物服务佣金（寄养/美容预约 15%）+ 品牌合作\n- **Phase 3 (12月+)**: 会员订阅（¥15/月）+ 线下活动票务\n\n**目标 Unit Economics**：ARPU ¥8/月（M12），CAC < ¥20\n\n## 8. 关键指标\n\n| 指标 | M3 | M6 | M12 |\n|------|----|----|-----|\n| 注册宠物数 | 3,000 | 15,000 | 80,000 |\n| 配对成功率 | 15% | 25% | 35% |\n| 线下见面转化 | 5% | 12% | 20% |\n| 月留存 | 30% | 40% | 50% |', '[{"step":"positioning","decision":"宠物第一人称社交","reason":"差异化竞争，避免做成又一个宠物版小红书"},{"step":"cold_start","decision":"线下活动驱动","reason":"社交产品需要密度，线下场景转化率高于纯线上"}]', 'submitted', 38, NOW() - INTERVAL '29 days', NOW() - INTERVAL '31 days', NOW() - INTERVAL '29 days'),

(1, 7, E'# Pawsitive — 宠物社交 MVP 方案\n\n## 1. 产品定位\n\n面向一线城市年轻养猫群体的轻社交工具。**先做猫，不做狗**——猫主人更宅、更依赖线上社交、分享欲更强。\n\n## 2. 市场切入\n\n中国宠物猫数量（2024）：7043 万只，首次超过犬（6844 万只）。养猫人群画像更年轻（25-30 岁占 52%）、更活跃于社交媒体。\n\n> 猫主人在小红书的日均发帖量是犬主人的 1.7 倍\n\n## 3. 功能设计\n\n### P0 — MVP\n- 猫咪主页（品种 / 性格 / 日常照片墙）\n- 附近猫友地图（3km 范围内的猫主人）\n- 猫咪日记（每日一图 + 一句话，轻量化记录）\n- 私信聊天\n\n### P1\n- 猫咪健康档案（疫苗、绝育、体重曲线）\n- 附近宠物医院 + 评价\n- 猫咪性格测试（问卷 + AI 分析）\n\n## 4. 技术选型\n- 微信小程序 MVP（降低获客门槛）\n- 后端 Go + PostgreSQL\n- 七牛云 CDN 图片存储\n\n## 5. 运营策略\n\n以城市为单位逐步覆盖：上海 → 北京 → 杭州 → 成都。每个城市先做到 500 活跃用户再开下一城。\n\n## 6. 指标\n\n| 指标 | 目标 |\n|------|------|\n| 首月注册 | 1,000 |\n| 日活/月活 | >20% |\n| 7 日留存 | 35% |', '[{"step":"scope","decision":"只做猫市场","reason":"猫主人线上社交需求更强，更容易冷启动"}]', 'submitted', 29, NOW() - INTERVAL '28 days', NOW() - INTERVAL '30 days', NOW() - INTERVAL '28 days'),

(1, 10, E'# 毛球星球 — 宠物版大众点评\n\n## 产品思路\n\n不做社交，做**宠物服务评价平台**。宠物社交太难冷启动，但"哪家宠物医院好"是刚需且搜索频次高。\n\n## 核心功能\n1. **服务点评**：宠物医院/美容/寄养的评分与评价\n2. **价格透明**：各服务商的价格对比\n3. **内容社区**：养宠经验分享（图文 + 短视频）\n4. **二手市场**：宠物用品闲置交易\n\n## 竞品分析\n- 大众点评有宠物类目但数据稀疏\n- 波奇宠物侧重电商，评价体系弱\n- **空白点**：没有人把宠物服务评价做透\n\n## 商业模式\n- 商家入驻费（基础免费，推广收费）\n- 预约佣金 10%\n- 广告收入\n\n## 技术\n- React Native 客户端\n- Node.js + MongoDB\n- 高德地图 API', '[]', 'submitted', 21, NOW() - INTERVAL '27 days', NOW() - INTERVAL '29 days', NOW() - INTERVAL '27 days'),

(1, 14, E'# WoofWoof — 宠物生活服务入口\n\n## 切入点\n\n围绕「遛狗」这一高频场景构建产品。遛狗是刚需、高频、有社交属性的行为，是最好的冷启动场景。\n\n## MVP 功能\n1. **遛弯地图**：标记附近适合遛狗的公园/绿地\n2. **遛弯组队**：发起遛弯邀约，附近狗友可加入\n3. **GPS 轨迹**：记录遛弯路线和时长\n\n## 扩展方向\n- 遛狗频率 → 健康数据\n- 遛狗路线 → 附近服务推荐\n- 遛狗社交 → 宠物用品团购\n\n## 技术架构\n- React Native 客户端\n- Go + gRPC 微服务后端\n- Redis GEO + ElasticSearch\n- 阿里云 OSS\n\n## 风险\n- 天气依赖强（雨天无人遛狗）\n- 地域密度要求高', '[]', 'submitted', 17, NOW() - INTERVAL '26 days', NOW() - INTERVAL '28 days', NOW() - INTERVAL '26 days');

-- Idea 2: 独立开发者收入仪表盘 (4 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(2, 1, E'# IndieMetrics — 独立开发者多平台收入聚合\n\n## 1. 问题验证\n\n我们在 Indie Hackers、V2EX 和 Twitter 对 50 位独立开发者做了调研：\n- 87% 使用 Excel/Notion 手动汇总多平台收入\n- 平均每月花 2.5 小时在财务数据整理上\n- 73% 希望有实时的统一收入视图\n- 最常用的收入平台：Stripe (82%)、Gumroad (41%)、App Store (35%)、LemonSqueezy (28%)\n\n## 2. 竞品分析\n\n| 产品 | 价格 | 优势 | 不足 |\n|------|------|------|------|\n| Baremetrics | $129+/月 | 功能全面，SaaS 指标丰富 | 贵，面向 SaaS 团队非个人 |\n| ChartMogul | $100+/月 | 数据准确，API 强 | 对接复杂，学习曲线陡 |\n| Indie Dashboard (开源) | 免费 | 可自托管 | 需自己部署维护，功能简陋 |\n\n**市场缺口**：没有专为独立开发者设计的、轻量且价格友好的多平台收入看板。\n\n## 3. 核心功能\n\n### P0 — MVP\n| 功能 | 描述 |\n|------|------|\n| 多平台 OAuth 对接 | Stripe、Gumroad、LemonSqueezy、Paddle、App Store Connect |\n| 实时仪表盘 | 总收入、MRR、ARR、增长趋势图 |\n| 货币统一 | 30+ 货币自动换算，支持设定基准货币 |\n| 每日收入邮件 | 每天早上 8 点发送昨日收入快报 |\n\n### P1\n- 收入预测（基于历史数据的线性/季节性趋势）\n- 目标追踪（设定月/年收入目标，实时进度条）\n- 公开收入页面（可选，Build in Public 场景）\n- 收入日历热力图\n\n### P2\n- 税务报表导出（按国家/地区）\n- Webhook 集成（新订单/退款通知到 Slack/Discord）\n- 团队视图（多人查看）\n\n## 4. 技术方案\n\n```\n[Next.js 14 App Router]\n        ↓\n[API Routes + Server Actions]\n        ↓\n[PostgreSQL (Supabase)] ← [Cron: 每 15min 同步各平台数据]\n        ↓\n[Vercel 部署] + [Resend 邮件服务]\n```\n\n- **全栈 Next.js**：一个人能 hold 住，部署在 Vercel\n- **数据同步**：Vercel Cron Job 每 15 分钟拉取各平台 API\n- **OAuth 2.0**：各平台标准接入，token 加密存储\n- **图表**：Recharts，SSR 友好\n\n## 5. 数据模型\n\n```sql\nusers (id, email, base_currency, created_at)\nconnections (id, user_id, platform, access_token_enc, status)\ntransactions (id, connection_id, amount, currency, type, occurred_at)\ndaily_snapshots (id, user_id, date, total_revenue, mrr)\n```\n\n## 6. 定价\n\n| 档位 | 价格 | 功能 |\n|------|------|------|\n| Free | $0 | 3 数据源，30 天历史 |\n| Pro | $9/月 | 无限数据源，完整历史，邮件报告，导出 |\n| Team | $29/月 | Pro + 多人查看 + API 接口 |\n\n**目标**：M6 达到 500 Pro 用户 = $4,500 MRR\n\n## 7. 增长策略\n\n1. 产品自身 Build in Public（每周公开 IndieMetrics 自己的收入数据）\n2. 公开收入页面作为自然传播渠道\n3. Indie Hackers / Twitter 社区运营\n4. ProductHunt Launch Day\n\n## 8. 风险\n\n| 风险 | 应对 |\n|------|------|\n| 平台 API 变更/限流 | 本地缓存 + 降级策略，保证核心数据不丢 |\n| 单人开发瓶颈 | MVP 极简，先做 Stripe + Gumroad 两个平台 |\n| 数据安全顾虑 | SOC2 合规路线，token 端到端加密 |', '[{"step":"market_research","decision":"聚焦独立开发者而非 SaaS 团队","reason":"竞品都在做 SaaS 团队市场，个人开发者被忽视且价格敏感"},{"step":"mvp_scope","decision":"首版只做 Stripe + Gumroad","reason":"覆盖 80% 用户的主要收入源，降低开发量"},{"step":"pricing","decision":"Free 档保留功能限制","reason":"转化漏斗必须有免费入口"}]', 'submitted', 55, NOW() - INTERVAL '27 days', NOW() - INTERVAL '28 days', NOW() - INTERVAL '27 days'),

(2, 6, E'# RevenuePulse — 极简收入追踪\n\n## 1. 产品理念\n\n**只做一件事：打开就能看到今天赚了多少。**\n\n不做财务分析工具，不做税务助手。只做一个最美最快的收入看板，像看天气 App 一样看收入。\n\n## 2. 差异化\n\n| 维度 | 竞品（Baremetrics等） | RevenuePulse |\n|------|----------------------|-------------|\n| 设置时间 | 30min+ | 5min |\n| 学习成本 | 需理解 SaaS 指标 | 零学习 |\n| 价格 | $100+/月 | $5/月 |\n| 目标用户 | SaaS 团队 | 个人开发者 |\n\n## 3. MVP 功能（2 周交付）\n\n只做 3 个页面：\n1. **Today**: 今日收入大数字 + 与昨日对比\n2. **Trend**: 30 天收入折线图\n3. **Sources**: 各平台收入占比饼图\n\n数据源第一版只接 Stripe + Gumroad（覆盖 80% 独立开发者）。\n\n## 4. 技术选型\n\n- **Next.js 14** App Router 全栈\n- **Vercel** 部署（零运维）\n- **Supabase** PostgreSQL + Auth\n- **Vercel Cron** 每小时同步数据\n- OAuth 接入 Stripe Connect 和 Gumroad API\n\n一个人两周可以做完 MVP。\n\n## 5. 商业模式\n\n极简定价：\n- Free: 1 个数据源，7 天历史\n- Pro: $5/月，无限数据源，完整历史\n\n**目标 M6**: 1000 Pro 用户 = $5,000 MRR\n\n## 6. 冷启动\n\n在 Twitter 做 Build in Public，每天晒自己产品的收入截图（用 RevenuePulse 截图），形成"产品即营销"。', '[{"step":"scope","decision":"极简主义","reason":"个人项目必须控制复杂度，两周内必须能上线"}]', 'submitted', 33, NOW() - INTERVAL '26 days', NOW() - INTERVAL '27 days', NOW() - INTERVAL '26 days'),

(2, 9, E'# MoneyBoard — 移动优先的收入看板\n\n## 1. 竞品调研\n\n| 产品 | 价格 | 移动端 | 独立开发者友好 |\n|------|------|--------|---------------|\n| Baremetrics | $129/月 | 有但体验差 | 否 |\n| ChartMogul | $100/月 | 无 | 否 |\n| Indie Dashboard | 免费 | 无 | 是但需自部署 |\n| **MoneyBoard** | $5/月 | 核心 | 是 |\n\n## 2. 核心洞察\n\n独立开发者最常看收入的场景是**早上起床刷手机**和**晚上睡前**。桌面端看板的使用频率远低于移动端。\n\n## 3. 产品设计\n\n### 移动端核心体验\n- **Widget**：iOS/Android 桌面小组件，不打开 App 就能看到今日收入\n- **推送通知**：新订单实时推送，收入里程碑提醒\n- **收入日历**：热力图展示每天收入，一目了然\n\n### 独特功能\n- **里程碑系统**：首次 $100、首次 $1000……自动记录并可分享\n- **公开收入页**：可选公开，用于 Build in Public\n- **每日邮件**：每天一封收入快报\n\n## 4. 技术方案\n\n- React Native（移动端 MVP）\n- Next.js（公开收入页 + Web 管理后台）\n- Supabase（DB + Auth + Realtime）\n- Expo Push Notifications\n\n## 5. 定价\n\n**$5/月**，无免费档。理由：\n- 免费用户维护成本高但转化率低\n- $5 价格门槛极低，筛选出有付费意愿的用户\n- Build in Public 社区对付费工具接受度高\n\n## 6. 风险\n\n- 移动端开发成本高 → 先做 PWA 替代原生\n- 数据源有限 → MVP 只做 Stripe，后续按需求优先级加', '[{"step":"platform","decision":"移动端优先","reason":"核心使用场景是早起/睡前看手机"},{"step":"pricing","decision":"无免费档","reason":"低价筛选付费用户，降低运维成本"}]', 'submitted', 28, NOW() - INTERVAL '25 days', NOW() - INTERVAL '26 days', NOW() - INTERVAL '25 days'),

(2, 15, E'# 独立开发者收入仪表盘 — 工具 + 社区方案\n\n## 核心洞察\n\nBuild in Public 趋势下，越来越多独立开发者愿意公开收入数据。但目前公开方式很原始（Twitter 截图/Notion 手动更新）。\n\n## 产品 = 工具 + 社区\n\n**工具层**：多平台收入聚合看板（基本功能）\n**社区层**：匿名收入排行榜 + 里程碑广场 + 经验分享\n\n## 增长飞轮\n\n```\n用户连接收入数据 → 选择公开收入页面\n→ 分享到 Twitter/社区 → 其他开发者看到\n→ 好奇心驱动注册 → 更多公开页面\n→ 社区活跃度提升 → 媒体报道\n```\n\n## 差异化\n\n不是更好的看板工具，而是**独立开发者的收入社区**。\n\n## 功能\n\n1. 收入看板（基础功能，对标竞品）\n2. 公开收入页（vanity URL: moneyboard.dev/@username）\n3. 匿名排行榜（按 MRR 区间分组，保护隐私）\n4. 里程碑广场（自动发布"首次 $1000 MRR"等成就）\n\n## 商业模式\n\n- 工具免费，社区功能付费（$8/月）\n- 赞助：开发者工具品牌赞助排行榜\n- 招聘：企业可以在社区发布面向独立开发者的合作机会', '[]', 'submitted', 19, NOW() - INTERVAL '24 days', NOW() - INTERVAL '25 days', NOW() - INTERVAL '24 days');

-- Idea 3: AI 面试模拟器 (5 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(3, 3, E'# MockInterview AI — 24/7 AI 面试教练\n\n## 1. 产品愿景\n\n让每个求职者都有一个 24 小时随叫随到的面试教练。不是题库，不是模板——是一个能听你说话、给你实时反馈、记住你弱点的 AI 面试官。\n\n## 2. 市场分析\n\n- 全球面试准备市场规模约 $2.4B（Grand View Research 2024）\n- 中国年应届毕业生 1100+ 万，技术岗竞争比 15:1\n- 面试辅导时薪 ¥200-800，AI 方案可降至 ¥0.5/次\n\n## 3. 竞品对比\n\n| 产品 | 模式 | 价格 | 短板 |\n|------|------|------|------|\n| Pramp | 真人对练 | 免费 | 需要预约，质量不稳定 |\n| InterviewBit | 题库 + 视频 | $30/月 | 无互动，缺乏反馈 |\n| 牛客面试 | 题库 + 面经 | 部分付费 | 文字为主，无模拟 |\n| **MockInterview AI** | AI 实时模拟 | $19/月 | 全新赛道 |\n\n## 4. 核心功能\n\n### P0 — MVP\n| 功能 | 描述 | 技术 |\n|------|------|------|\n| 定制面试 | 选择公司+职位 → AI 生成针对性问题 | GPT-4 + 公司题库 |\n| 语音对话 | 真实面试体验，AI 扮演面试官 | Whisper STT + TTS |\n| 实时反馈 | STAR 法则评分 + 改进建议 | 自定义评估 prompt |\n| 弱项追踪 | 记录历史表现，自动推荐针对性练习 | 用户画像模型 |\n\n### P1\n- 视频面试（摄像头 + 表情/姿态分析）\n- 代码面试（集成在线编辑器 + AI 评审）\n- 英文面试模式\n- 面试复盘报告（可分享 PDF）\n\n### P2\n- 企业定制题库（B2B，企业上传自己的面试题）\n- 模拟群面/Case Study\n- 薪资谈判模拟\n\n## 5. 技术架构\n\n```\n[React Web App] ←WebSocket→ [Interview Engine]\n                                    ↓\n                    [OpenAI GPT-4] + [Whisper API]\n                                    ↓\n                    [评估引擎: STAR 评分 + 关键词匹配]\n                                    ↓\n                    [PostgreSQL] + [Redis Session]\n```\n\n- **语音链路**: 浏览器 MediaRecorder → WebSocket → Whisper STT → GPT-4 → TTS → 播放\n- **延迟优化**: Whisper streaming + GPT-4 streaming，目标端到端 < 2s\n- **题库**: 初期人工整理 500 题（覆盖 Top 50 公司），后续 AI 生成扩充\n\n## 6. 用户画像\n\n**Primary: 小张**\n- 25 岁，前端开发 2 年经验，准备跳槽大厂\n- 技术能力不错，但面试表达差（"知道但说不清楚"）\n- 每天下班后练习 30 分钟\n- 付费意愿：¥100/月内\n\n## 7. 商业模式\n\n| 档位 | 价格 | 功能 |\n|------|------|------|\n| Free | $0 | 每天 3 题文字模拟 |\n| Pro | $19/月 | 无限语音模拟 + 评估报告 + 弱项训练 |\n| Enterprise | 按需 | 企业定制题库 + 批量账号 |\n\n## 8. 关键指标\n\n| KPI | M3 | M6 |\n|-----|-----|-----|\n| 注册用户 | 5,000 | 20,000 |\n| Pro 转化率 | 5% | 8% |\n| 日均面试次数 | 500 | 3,000 |\n| NPS | 40 | 55 |', '[{"step":"market_research","decision":"技术面试+产品面试双线","reason":"技术岗最大但产品岗付费意愿更高"},{"step":"tech","decision":"先做语音后做视频","reason":"语音已足够模拟真实面试，视频增加复杂度但提升有限"},{"step":"pricing","decision":"$19/月","reason":"低于真人辅导但高于题库，卡在\"一杯咖啡钱\"心理锚点"}]', 'submitted', 61, NOW() - INTERVAL '25 days', NOW() - INTERVAL '26 days', NOW() - INTERVAL '25 days'),

(3, 8, E'# InterviewGPT — 端到端面试准备平台\n\n## 1. 核心差异\n\n不只是面试模拟器，而是**完整的面试准备旅程**：\n\n```\n简历分析 → 弱点定位 → 知识补充 → 定制训练 → 模拟面试 → 复盘改进\n```\n\n现有产品都在做"模拟面试"这一个环节，但求职者的痛点分布在整个链路上。\n\n## 2. 市场机会\n\n面试准备的核心用户分两类：\n- **应届生**（量大、付费弱）：主要需要系统化训练\n- **跳槽程序员**（量中、付费强）：主要需要针对性准备\n\n我们选择先做**跳槽程序员**——ARPU 更高，需求更明确。\n\n## 3. 功能设计\n\n### 旅程设计\n1. **Upload**: 上传简历 → AI 分析技能图谱和弱点\n2. **Plan**: 基于目标公司生成个性化准备计划（2-4 周）\n3. **Learn**: 针对弱点的知识卡片和刻意练习\n4. **Mock**: AI 模拟面试（文字 → 语音分阶段上线）\n5. **Review**: 面试后 AI 生成详细复盘报告\n\n### 技术亮点\n- 简历解析：OCR + NLP 提取技能点，构建个人技能图谱\n- 知识图谱：面试知识点之间的关联图，智能推荐学习路径\n- 模拟面试：GPT-4 扮演不同风格面试官（nice型/压力型/追问型）\n\n## 4. 分阶段交付\n\n| Phase | 功能 | Timeline |\n|-------|------|----------|\n| Phase 1 | 简历分析 + 文字模拟 + 复盘报告 | M1-M2 |\n| Phase 2 | 语音面试 + 个性化计划 | M3-M4 |\n| Phase 3 | 视频面试 + 肢体语言分析 + B2B | M5-M8 |\n\n## 5. 技术栈\n\n- 前端: Next.js 14 + Tailwind\n- 后端: Python FastAPI（AI pipeline）+ Go（业务逻辑）\n- AI: GPT-4 Turbo + Whisper + 自训练评分模型\n- DB: PostgreSQL + Qdrant（知识向量搜索）\n\n## 6. 商业模式\n\n- B2C: $15/月（个人版）\n- B2B: 按座位收费（培训机构/高校就业中心）\n- 估算 M12 ARR: $180K', '[{"step":"positioning","decision":"做全链路而非单点模拟","reason":"单点竞争易被复制，全链路体验形成壁垒"},{"step":"target","decision":"先做跳槽程序员","reason":"付费能力强，需求明确，获客渠道清晰"}]', 'submitted', 48, NOW() - INTERVAL '24 days', NOW() - INTERVAL '25 days', NOW() - INTERVAL '24 days'),

(3, 11, E'# 面试达人 — 中国技术面试专属平台\n\n## 1. 精准定位\n\n**只做技术面试，只做中国市场。**\n\n为什么？中国大厂面试有其独特性：八股文、手撕代码、系统设计三板斧，海外产品不理解这些。\n\n## 2. 题库来源与壁垒\n\n| 来源 | 数量 | 获取方式 |\n|------|------|----------|\n| 牛客网面经 | 5,000+ | 公开爬取 + 结构化 |\n| 力扣讨论区 | 3,000+ | 按公司标签聚合 |\n| 大厂内推群 | 1,000+ | 社区用户贡献（脱敏处理） |\n| AI 生成 | 无限 | 基于真题风格生成变体 |\n\n## 3. 核心功能\n\n1. **按公司/岗位/轮次筛选面试题**\n2. **AI 模拟面试**：八股问答 + 手撕代码 + 系统设计三种模式\n3. **代码面试**：集成 Monaco Editor，AI 实时评审代码\n4. **面试官风格**：可选择「友善型」「压力型」「追问型」\n5. **社区互评**：用户互相模拟面试、互相评价\n\n## 4. 技术方案\n\n- 前端: React + Monaco Editor + WebSocket\n- 后端: Go + Echo\n- AI: Claude API（中文理解更好）+ 自训练八股文评分模型\n- 代码执行: 沙箱化 Docker 容器\n\n## 5. 商业模式\n\n- 免费：每天 5 题文字练习\n- VIP ¥39/月：语音模拟 + 代码面试 + 详细报告\n- 年卡 ¥299：适合应届生（秋招周期 4 个月）\n\n## 6. 风险\n\n- 牛客网可能做类似功能 → 差异化在 AI 模拟体验\n- 八股文面试趋势在减弱 → 增加系统设计和场景题', '[{"step":"market","decision":"只做中国市场","reason":"中国面试文化独特，海外产品水土不服"}]', 'submitted', 35, NOW() - INTERVAL '23 days', NOW() - INTERVAL '24 days', NOW() - INTERVAL '23 days'),

(3, 16, E'# PrepAI — 游戏化面试训练\n\n## 核心思路\n\n把面试准备做成 Duolingo——每天 5 分钟，持续进步。\n\n## 游戏化机制\n\n- **每日挑战**: 3 道随机面试题，限时作答\n- **连续打卡**: 连续 7 天送 Pro 体验券\n- **段位系统**: 青铜 → 白银 → 黄金 → 钻石（按答题正确率）\n- **排行榜**: 按职位/公司分组，看你在同目标求职者中的水平\n- **成就徽章**: "连续 30 天""击败 90% 用户""系统设计满分"\n\n## 功能\n\n1. 每日挑战（3 题，5 分钟）\n2. 专项训练（按知识点分类：算法/系统设计/行为面试）\n3. AI 即时评分 + 参考答案\n4. 学习路径推荐\n\n## 技术\n\n- React + Tailwind CSS 前端\n- Python FastAPI 后端\n- OpenAI API（题目生成 + 答案评估）\n- Redis（排行榜 + 连续天数缓存）\n\n## 定价\n\n- Free: 每日挑战 + 基础评分\n- Pro $9/月: 专项训练 + 详细评分 + 排行榜', '[]', 'submitted', 26, NOW() - INTERVAL '22 days', NOW() - INTERVAL '24 days', NOW() - INTERVAL '22 days'),

(3, 20, E'# HireReady — 面试策略教练\n\n## 独特视角\n\n不做面试模拟，做**面试策略**。帮你在面试前做好 80% 的功课。\n\n## 核心功能\n\n1. **JD 深度解析**: 上传 JD → AI 提取关键技能要求 + 隐藏需求\n2. **公司调研报告**: 自动汇总公司技术博客、GitHub、面经信息\n3. **策略建议**: 每轮面试该准备什么、该强调什么\n4. **薪资谈判助手**: 基于市场数据给出谈判策略和话术\n\n## 目标用户\n\n有 3-5 年经验、准备跳槽到更好公司的中级工程师。他们不缺技术能力，缺的是"面试软实力"。\n\n## 技术\n\n- Web App (Next.js)\n- AI: GPT-4 + 爬虫聚合公司信息\n- 数据源: 脉脉/Glassdoor/GitHub/公司官网', '[]', 'submitted', 18, NOW() - INTERVAL '21 days', NOW() - INTERVAL '23 days', NOW() - INTERVAL '21 days');

-- Idea 4: 团队知识图谱 (4 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(4, 5, E'# TeamBrain — 从工具链自动构建团队知识图谱\n\n## 1. 问题定义\n\n> 新人问："这个服务是谁写的？出了问题找谁？"\n> 答："你去问 Tom，他可能知道……"\n\n这种场景在 10-50 人的技术团队中每天发生数十次。知识分散在 Git、Slack、Notion 和个人脑中，没有统一的索引。\n\n**量化痛点**：\n- 新人 onboarding 平均 3-6 个月才能独立定位问题\n- 工程师每周花 4.2 小时在"找人问"上（Atlassian 2024 调研）\n- 关键人离职时带走大量隐性知识\n\n## 2. 解决方案\n\n自动从团队已有工具链（Git、Slack、Notion、Jira）提取数据，构建 **人 ↔ 知识 ↔ 项目** 的三维关联图谱。\n\n### 核心查询场景\n| 查询 | 结果 |\n|------|------|\n| "谁最了解 payment 模块？" | Git 贡献者排名 + 相关 PR 链接 |\n| "Redis 缓存策略的讨论在哪？" | Slack 消息链接 + Notion 文档 |\n| "上次数据库迁移出了什么问题？" | Postmortem 文档 + 相关 Jira ticket |\n| "这个 API 最近谁改过？" | Git blame 增强视图 |\n\n## 3. 竞品分析\n\n| 产品 | 定位 | 价格 | 短板 |\n|------|------|------|------|\n| Guru | 知识库 wiki | $10/人/月 | 需要手动录入，不自动 |\n| Tettra | 内部 wiki | $8/人/月 | 无代码上下文理解 |\n| Glean | 企业搜索 | $15/人/月 | 面向大企业，中小团队用不起 |\n\n**我们的差异**：零手动输入，全自动从已有工具链提取。\n\n## 4. 技术架构\n\n```\n[Data Connectors]\n  ├── GitHub (commits, PRs, code review)\n  ├── Slack (conversations, threads)\n  ├── Notion (documents, wikis)\n  └── Jira (tickets, comments)\n        ↓ (增量同步)\n[ETL Pipeline: Go workers]\n        ↓\n[知识图谱引擎]\n  ├── PostgreSQL (关系数据)\n  ├── Qdrant (向量嵌入)\n  └── Neo4j (图查询)\n        ↓\n[查询层: GraphQL API]\n        ↓\n[Slack Bot] + [Web Dashboard] + [VS Code Extension]\n```\n\n- **增量索引**：Webhook + 定时同步，不影响日常工作流\n- **NLP**: OpenAI Embedding + GPT-4 理解自然语言查询\n- **权限继承**：继承源系统权限，Git private repo 的知识仅相关人可见\n\n## 5. MVP 路线图\n\n| Phase | 范围 | Timeline |\n|-------|------|----------|\n| MVP | GitHub connector only → "谁了解 X 模块" | M1-M2 |\n| V1 | + Slack connector → 全文对话搜索 | M3-M4 |\n| V2 | + Notion → 完整知识图谱 + Web UI | M5-M6 |\n| V3 | VS Code 插件 + 自动 reviewer 推荐 | M7-M8 |\n\n## 6. 商业模式\n\n- 开源核心（GitHub connector），商业版增加 Slack/Notion/Jira\n- 定价：$6/人/月（10 人团队 = $60/月）\n- 私有化部署：$500/月起\n\n## 7. 成功指标\n\n| KPI | M3 | M6 |\n|-----|-----|-----|\n| 接入团队数 | 20 | 100 |\n| 日均查询数 | 50 | 500 |\n| 新人 onboarding 时间缩短 | 20% | 40% |', '[{"step":"scope","decision":"MVP 只做 GitHub connector","reason":"Git 数据结构化程度最高，提取难度最低"},{"step":"architecture","decision":"Neo4j 图数据库","reason":"知识关联天然是图结构，图查询性能远优于关系型"}]', 'submitted', 44, NOW() - INTERVAL '23 days', NOW() - INTERVAL '24 days', NOW() - INTERVAL '23 days'),

(4, 12, E'# KnowGraph — 最小可用的 "Who Knows What"\n\n## 1. 核心主张\n\n不做大知识图谱，先做最小有价值的场景：**谁了解什么（Who Knows What）**。\n\n## 2. MVP：只对接 GitHub\n\n通过 commit 历史、PR review 记录和 code ownership 数据，自动建立「人 → 模块 → 技能」的映射。\n\n**数据提取规则**：\n| 数据源 | 提取内容 | 权重 |\n|--------|----------|------|\n| Git commits | 文件路径 → 模块归属 | 40% |\n| PR reviews | 审核模块 → 熟悉度 | 30% |\n| Issue comments | 讨论主题 → 知识标签 | 20% |\n| CODEOWNERS | 显式声明的负责人 | 10% |\n\n## 3. 交互形式\n\n**Slack Bot**（主要入口）：\n```\n@knowgraph 谁了解 payment 模块？\n→ 1. @alice (85 commits, 23 PRs reviewed)\n   2. @bob (42 commits, 15 PRs reviewed)\n   3. @charlie (12 commits, 8 PRs reviewed)\n```\n\n**Web Dashboard**（辅助）：\n- 团队知识分布热力图\n- 模块负责人矩阵\n- Bus Factor 风险预警（某模块只有 1 人了解）\n\n## 4. 技术选型\n\n- GitHub App（OAuth + Webhook）\n- Go 后端 + PostgreSQL\n- Slack Bot API\n- 无需向量数据库（MVP 阶段基于关键词匹配足够）\n\n## 5. 定价\n\n- Free: 5 人以下团队\n- Pro: $4/人/月\n- 估算: 50 个付费团队（平均 15 人）= $3,000 MRR\n\n## 6. 冷启动\n\n在 GitHub Marketplace 上架，利用 GitHub 生态获取自然流量。', '[{"step":"scope","decision":"只做 Who Knows What","reason":"一个场景做透比十个场景做浅更有价值"},{"step":"channel","decision":"Slack Bot 为主入口","reason":"工程师日常在 Slack，不会主动打开新工具"}]', 'submitted', 37, NOW() - INTERVAL '22 days', NOW() - INTERVAL '23 days', NOW() - INTERVAL '22 days'),

(4, 17, E'# 团队大脑 — 隐私优先的知识引擎\n\n## 核心差异：隐私优先\n\n企业知识数据极其敏感。我们的方案：**所有数据处理在客户自己的基础设施内完成**。\n\n## 技术方案\n\n- **向量数据库**: Qdrant（自托管）存储知识嵌入\n- **LLM**: 支持 self-hosted LLM（Ollama/vLLM）或 API（OpenAI/Claude）\n- **增量索引**: 每小时同步新数据，不影响日常工作流\n- **部署方式**: Docker Compose 一键部署到客户 VPC\n\n## 隐私设计\n\n| 层级 | 策略 |\n|------|------|\n| 数据存储 | 全部本地，不出客户网络 |\n| LLM 调用 | 支持本地模型，可选 API |\n| 权限控制 | 继承 GitHub/Slack 权限 |\n| 审计日志 | 所有查询可追溯 |\n\n## 定位\n\n面向对数据安全有严格要求的团队（金融、医疗、政府外包）。\n\n## 定价\n\n- 开源版：基础功能免费\n- Enterprise: $800/月（含技术支持 + SLA）', '[]', 'submitted', 22, NOW() - INTERVAL '21 days', NOW() - INTERVAL '22 days', NOW() - INTERVAL '21 days'),

(4, 21, E'# CollectiveIQ — 主动式知识推送\n\n## 差异化\n\n不是搜索工具，是**知识发现**工具。不需要你主动提问——它会在你需要的时候，主动推送你可能需要但不知道存在的知识。\n\n## 场景\n\n1. **写代码时**：检测到你在实现某功能 → 推送团队内已有的类似实现\n2. **提 PR 时**：分析改动范围 → 自动推荐最合适的 reviewer\n3. **新人入职**：根据你被分配的任务 → 生成个性化学习路径\n4. **Postmortem 时**：搜索历史上类似的事故和解决方案\n\n## 技术\n\n- VS Code 插件（主要入口）\n- 代码语义分析 + 向量相似度匹配\n- GitHub + Slack + Notion 数据源\n\n## 风险\n\n- 推送干扰（过度通知）→ 需要精细的相关度阈值\n- 数据隐私 → 敏感代码仓库需要特殊处理', '[]', 'submitted', 15, NOW() - INTERVAL '20 days', NOW() - INTERVAL '21 days', NOW() - INTERVAL '20 days');

-- Idea 5: 播客笔记助手 (4 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(5, 1, E'# PodNotes — 播客内容结构化引擎\n\n## 1. 痛点验证\n\n> "上周某期播客里有人推荐了一本书，但我完全想不起来是哪期了。"\n\n我们采访了 30 位重度播客听众（周听 5h+）：\n- 92% 遇到过"想回顾某个观点但找不到"的情况\n- 67% 曾手动做过播客笔记但坚持不下来\n- 平均每人订阅 12 个播客，但只有 3 个能记住核心内容\n\n## 2. 核心功能\n\n### P0 — MVP\n| 功能 | 描述 | 技术 |\n|------|------|------|\n| 自动转录 | 音频 → 带时间戳的文字 | Whisper Large V3 |\n| 结构化摘要 | 每期生成：核心观点、提到的人/书/工具 | GPT-4 Turbo |\n| 关键信息卡片 | 自动提取书名/工具名/观点，归类索引 | NER + 分类模型 |\n| 全文搜索 | 搜索所有听过的播客的任何内容 | Meilisearch |\n| 语义搜索 | "那期讲创业融资的播客" | OpenAI Embedding + Qdrant |\n\n### P1\n- 时间戳精确定位（点击摘要跳到对应播放位置）\n- 跨集知识图谱（不同播客提到的相同主题自动关联）\n- 笔记导出（Notion / Obsidian / Markdown）\n- 每周知识回顾邮件\n\n## 3. 技术架构\n\n```\n[用户提交播客 URL / RSS]\n        ↓\n[音频下载 + 预处理]\n        ↓\n[Whisper STT → 带时间戳转录]\n        ↓\n[GPT-4 结构化提取]\n  ├── 摘要\n  ├── 关键实体（人/书/工具）\n  └── 核心观点列表\n        ↓\n[PostgreSQL + Qdrant + Meilisearch]\n        ↓\n[Web App (Next.js)]\n```\n\n处理一期 1 小时播客的成本约 $0.15（Whisper $0.006/min + GPT-4 ~$0.08）。\n\n## 4. 竞品分析\n\n| 产品 | 功能 | 不足 |\n|------|------|------|\n| Snipd | 片段标记 + AI 摘要 | 不支持中文播客 |\n| Podwise | 中文摘要 | 搜索弱，无知识关联 |\n| Readwise Reader | 阅读标注 | 播客支持初级 |\n\n## 5. 商业模式\n\n- Free: 每月 5 期播客处理\n- Pro $8/月: 无限处理 + 语义搜索 + 导出\n- Team $20/月: 团队共享知识库\n\n## 6. 指标\n\n| KPI | M3 | M6 |\n|-----|-----|-----|\n| 处理播客集数 | 5,000 | 30,000 |\n| MAU | 1,000 | 5,000 |\n| Pro 转化 | 8% | 12% |', '[{"step":"tech","decision":"Whisper + GPT-4 pipeline","reason":"Whisper 中文识别率 95%+，GPT-4 结构化提取最可靠"},{"step":"scope","decision":"先做 Web，不做浏览器插件","reason":"MVP 速度优先，插件审核周期长"}]', 'submitted', 39, NOW() - INTERVAL '21 days', NOW() - INTERVAL '22 days', NOW() - INTERVAL '21 days'),

(5, 13, E'# 播客回忆录 — 播客平台的笔记叠加层\n\n## 1. 产品形态\n\n**浏览器插件 + 移动 App**，在你常用的播客平台（小宇宙/Apple Podcasts/Spotify）上叠加一层智能笔记。\n\n**设计理念**：不改变用户听播客的习惯，而是在现有习惯上增加"记忆"。\n\n## 2. 核心功能\n\n1. **一键标记**：收听时按一下 = 标记精彩片段（保留前后 30s 上下文）\n2. **AI 总结**：每期播客自动生成 3 个核心观点\n3. **知识卡片**：自动提取人名、书名、工具名，归类到个人知识库\n4. **跨集搜索**：在所有听过的播客中搜索任何内容\n\n## 3. 用户场景\n\n| 场景 | 功能 |\n|------|------|\n| 通勤时听到好观点 | 一键标记 |\n| 想回忆上周听到的书名 | 搜索知识卡片 |\n| 周末写文章需要引用 | 搜索 → 跳转到精确时间戳 |\n\n## 4. 技术\n\n- 浏览器插件 (Chrome/Firefox): 注入播放器控件\n- React Native App\n- Whisper API + GPT-4 后端处理\n- 插件与 App 账号同步\n\n## 5. 商业模式\n\n- Free: 每月 10 次标记，基础搜索\n- Premium ¥15/月: 无限标记 + AI 总结 + 知识卡片 + 导出', '[{"step":"form","decision":"插件优先","reason":"不让用户换播客 App，降低使用门槛"}]', 'submitted', 30, NOW() - INTERVAL '20 days', NOW() - INTERVAL '21 days', NOW() - INTERVAL '20 days'),

(5, 18, E'# EarMark — 播客内容的搜索引擎\n\n## 独特定位\n\n不是笔记工具，是**播客世界的 Google**。可以用自然语言搜索任何播客曾经说过的话。\n\n## 技术壁垒\n\n1. **中英文播客索引**：自建爬虫，覆盖 10 万+节目、200 万+单集\n2. **全文转录**：Whisper 批量处理，构建最大的播客文字库\n3. **语义搜索**："那期讲 YC 创业的" → 精确匹配到具体集和时间戳\n4. **知识图谱**：不同播客提到的相同主题自动关联\n\n## 使用场景\n\n- 搜索"硅谷创业公司的融资策略" → 返回 20 个相关播客片段\n- 搜索"张一鸣管理方法" → 返回所有提到张一鸣管理理念的播客\n\n## 商业模式\n\n- Free: 搜索结果预览\n- Pro $12/月: 完整转录 + 时间戳跳转 + API\n- B2B: 播客平台/媒体公司 API 接入\n\n## 挑战\n\n- 音频处理成本高 → 冷数据用 Whisper small，热数据用 Whisper large\n- 版权问题 → 只索引公开 RSS，不存储原始音频', '[]', 'submitted', 25, NOW() - INTERVAL '19 days', NOW() - INTERVAL '20 days', NOW() - INTERVAL '19 days'),

(5, 22, E'# ListenLearn — 把播客变成学习工具\n\n## 核心目标\n\n把播客从「听了就忘的消遣」变成「可回顾可检索的学习资料」。\n\n## 功能\n\n1. **AI 思维导图**：每期播客自动生成思维导图\n2. **关联推荐**：听完一期后推荐讲同一主题的其他播客\n3. **笔记导出**：一键导出到 Notion / Obsidian / Logseq\n4. **学习追踪**：记录学习进度，生成月度学习报告\n5. **间隔重复**：重要观点定期推送复习\n\n## 技术\n\n- Web App (Next.js)\n- AI: Whisper + GPT-4 + 思维导图生成\n- Notion/Obsidian API 导出\n\n## 定价\n\n- Free: 3 期/月\n- Pro ¥12/月: 无限 + 导出 + 间隔重复', '[]', 'submitted', 16, NOW() - INTERVAL '18 days', NOW() - INTERVAL '19 days', NOW() - INTERVAL '18 days');

-- Idea 6: 远程团队异步站会 (3 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(6, 3, E'# AsyncStandup — 用异步取代每日站会\n\n## 1. 核心理念\n\n**站会的目的是同步信息，不是开会。** 跨时区团队每天凑齐所有人开会是反生产力的。\n\n## 2. 目标用户\n\n跨 3+ 时区的远程技术团队（10-30 人）。这类团队传统站会要么有人深夜参加，要么录像后没人看。\n\n## 3. 产品设计\n\n### 核心流程\n```\n每日下班前 → 花 2 分钟填写 3 个问题：\n1. 今天完成了什么？\n2. 明天计划做什么？\n3. 是否有阻塞？\n        ↓\nAI 自动生成团队日报\n        ↓\n发送到 Slack 频道 + 邮件\n```\n\n### 智能功能\n| 功能 | 描述 |\n|------|------|\n| AI 日报 | 汇总全队更新，提取关键信息 |\n| 异常检测 | 连续 3 天同一任务未完成 → 自动提醒 |\n| 阻塞匹配 | 有人被阻塞 → AI 推荐谁能帮忙 |\n| 周报生成 | 基于每日更新自动生成周报 |\n| 趋势分析 | 团队生产力趋势、阻塞率变化 |\n\n### Slack 深度集成\n- `/standup` 命令直接提交更新\n- 自动发送团队摘要到指定频道\n- Thread 回复进行异步讨论\n- Slack reminder 提醒未提交的成员\n\n## 4. 竞品\n\n| 产品 | 价格 | 短板 |\n|------|------|------|\n| Geekbot | $3/人/月 | 无 AI 分析，纯表单 |\n| Standuply | $2/人/月 | 体验老旧，集成弱 |\n| Range | $8/人/月 | 太重，不专注 |\n\n**我们的差异**: AI 驱动的异常检测 + 阻塞自动匹配，从"收集信息"进化到"解决问题"。\n\n## 5. 技术\n\n- Slack App (Bolt framework)\n- Go 后端 + PostgreSQL\n- OpenAI API (日报生成 + 异常检测)\n- Cron job (每日汇总 + 提醒)\n\n## 6. 定价\n\n- Free: 5 人以下\n- Pro $3/人/月: 无限人 + AI 功能 + 分析\n\n## 7. 成功指标\n\n| KPI | 目标 |\n|------|------|\n| 每日提交率 | >85% |\n| 节省会议时间 | 30min/天/团队 |\n| M6 付费团队 | 100 |', '[{"step":"form","decision":"Slack-first","reason":"工程师已经在 Slack，不需要打开新工具"},{"step":"differentiator","decision":"AI 异常检测","reason":"竞品只收集数据，我们还分析和行动"}]', 'submitted', 31, NOW() - INTERVAL '19 days', NOW() - INTERVAL '20 days', NOW() - INTERVAL '19 days'),

(6, 9, E'# DailySync — 短视频异步站会\n\n## 1. 核心差异\n\n**不只是文字，支持 60 秒短视频更新。** 有时候面对面说两句话比写一段文字更高效、更有人情味。\n\n远程团队最大的问题不是信息同步，是**情感连接缺失**。短视频让你看到队友的脸，听到他们的声音。\n\n## 2. 功能\n\n1. **多格式更新**：文字 / 60s 短视频 / 语音条，自由选择\n2. **AI 翻译**：视频/语音自动翻译字幕（解决跨语言团队）\n3. **任务联动**：与 Jira/Linear/GitHub 关联，自动填充"完成了什么"\n4. **周报一键生成**：AI 从每日更新汇总成周报\n5. **回顾模式**：快速浏览所有人今天的更新，像刷短视频一样\n\n## 3. 技术\n\n- React + React Native 跨平台\n- WebRTC 录制短视频\n- Whisper 转录 + GPT-4 翻译\n- Jira / Linear / GitHub OAuth 集成\n- 七牛云存储视频\n\n## 4. 定价\n\n- Free: 文字更新，5 人以下\n- Pro $5/人/月: 视频 + AI 翻译 + 任务联动\n\n## 5. 风险\n\n- 视频存储成本高 → 30 天自动归档\n- 用户不愿意录视频 → 文字作为 fallback', '[{"step":"differentiator","decision":"短视频模式","reason":"远程团队缺乏情感连接，视频比文字更有温度"}]', 'submitted', 24, NOW() - INTERVAL '18 days', NOW() - INTERVAL '19 days', NOW() - INTERVAL '18 days'),

(6, 19, E'# TeamPulse — 异步站会 + 团队健康度\n\n## 核心功能\n\n1. **异步签到**：每天填写进展 + 阻塞\n2. **情绪追踪**：今天工作状态 1-5 分（匿名）\n3. **阻塞匹配**：有人阻塞 → 系统推荐能帮忙的人\n4. **管理者仪表盘**：团队生产力 + 情绪趋势图\n\n## 差异化\n\n加入**情绪维度**。管理者不只看到任务进展，还能感知团队士气。连续低分预警 → 提醒 1:1 沟通。\n\n## 技术\n\n- Slack Bot + Web Dashboard\n- Go + PostgreSQL\n- 情绪趋势分析（简单统计模型）\n\n## 定价\n\n$4/人/月，10 人起', '[]', 'submitted', 14, NOW() - INTERVAL '17 days', NOW() - INTERVAL '18 days', NOW() - INTERVAL '17 days');

-- Idea 7: GitHub Star 追踪器 (4 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(7, 2, E'# StarKeeper — GitHub Star 智能管理平台\n\n## 1. 问题\n\n开发者平均 Star 了 200+ 个项目，但 GitHub 提供的管理功能几乎为零：无法搜索、无法分类、无法追踪更新。Star 列表变成了一个无法使用的垃圾堆。\n\n## 2. 核心功能\n\n| 功能 | 描述 | 价值 |\n|------|------|------|\n| 自动分类 | 根据 README、语言、topic 自动打标签 | 管理效率 |\n| 更新追踪 | 新版本/breaking change 通知 | 不错过重要更新 |\n| 死项目检测 | 6 月无更新标记为 archived | 清理收藏 |\n| 替代品推荐 | 项目不维护 → 推荐同类活跃项目 | 技术选型辅助 |\n| 搜索增强 | 全文搜索 Star 过的项目 | 快速找到 |\n| 项目健康度 | 综合评分（活跃度/社区/文档/issue响应） | 选型参考 |\n\n## 3. 技术方案\n\n```\n[GitHub OAuth] → [同步 Star 列表]\n        ↓\n[分类引擎: GPT-4 分析 README]\n        ↓\n[PostgreSQL + Meilisearch]\n        ↓\n[Cron: 每日检查 Release]\n        ↓\n[通知: Email + Slack webhook]\n```\n\n- GitHub API v4 (GraphQL) 获取 Star 数据\n- GPT-4 分析 README 自动生成分类标签\n- GitHub Webhook 监听 Release 事件\n- LLM 自动总结 changelog 关键变更\n\n## 4. 竞品\n\n| 产品 | 功能 | 价格 | 问题 |\n|------|------|------|------|\n| Astral | Star 管理 | 免费 | 已停止维护 |\n| Star History | Star 趋势图 | 免费 | 只看趋势，不管理 |\n| **StarKeeper** | 管理+追踪+推荐 | $5/月 | 全面方案 |\n\n## 5. 商业模式\n\n- Free: 100 Star 管理 + 基础分类\n- Pro $5/月: 无限 Star + 通知 + 健康度评分 + 替代品推荐\n\n## 6. 指标\n\n| KPI | M3 | M6 |\n|-----|-----|-----|\n| 注册用户 | 3,000 | 10,000 |\n| Pro 转化 | 5% | 8% |\n| MRR | $750 | $4,000 |', '[{"step":"scope","decision":"Web App 优先","reason":"GitHub 用户主要在桌面端工作"},{"step":"tech","decision":"GPT-4 自动分类","reason":"人工标签不可扩展，LLM 分析 README 准确率 >90%"}]', 'submitted', 46, NOW() - INTERVAL '17 days', NOW() - INTERVAL '18 days', NOW() - INTERVAL '17 days'),

(7, 8, E'# GitStar Pro — GitHub 增强浏览器插件\n\n## 1. 产品形态\n\n不做独立 App，做**Chrome/Firefox 插件**，直接增强 GitHub Star 页面。\n\n**设计原则**：用户不需要离开 GitHub，所有功能嵌入到现有页面中。\n\n## 2. 功能\n\n### 注入到 GitHub Star 页面\n- **搜索增强**：在 Star 列表顶部加搜索框（GitHub 原生没有！）\n- **筛选面板**：按语言/标签/更新时间筛选\n- **排序选项**：按 Star 数/最近更新/健康度排序\n- **自定义标签**：给每个项目打标签\n\n### 项目详情增强\n- **健康度评分**：综合 commit 频率/issue 响应/contributor 数\n- **Release 高亮**：新版本标记，changelog 摘要\n\n### 通知\n- 浏览器通知：关注项目有新 Release\n- 每周精选邮件：你 Star 的项目本周重要更新\n\n## 3. 技术\n\n- Chrome Extension Manifest V3\n- Content Script 注入 GitHub 页面\n- Background Worker 定时检查 Release\n- Supabase 后端（用户数据 + 标签同步）\n\n## 4. 商业模式\n\n- Free: 搜索 + 筛选（核心功能免费，做装机量）\n- Pro $3/月: 通知 + 健康度 + 每周邮件\n\n## 5. 冷启动\n\nChrome Web Store + ProductHunt + Hacker News。插件天然有口碑传播。', '[{"step":"form","decision":"浏览器插件","reason":"不需要用户改变习惯，GitHub 页面原地增强"}]', 'submitted', 34, NOW() - INTERVAL '16 days', NOW() - INTERVAL '17 days', NOW() - INTERVAL '16 days'),

(7, 14, E'# Constellation — Star 社交发现\n\n## 核心理念\n\n**看看和你 Star 品味相似的开发者都在关注什么。**\n\n你 Star 了 200 个项目，另一个人也 Star 了其中 150 个——他多 Star 的那 50 个，很可能正是你需要的。\n\n## 功能\n\n1. **Star 品味匹配**：找到和你 Star 重合度最高的开发者\n2. **推荐引擎**："和你口味相似的人还 Star 了这些"\n3. **协作列表**：类似 awesome-list 但动态更新、可协作\n4. **项目讨论**：围绕 Star 项目的轻量社区\n\n## 技术\n\n- 协同过滤算法（User-based CF）\n- GitHub OAuth + Star 数据\n- Next.js + PostgreSQL\n\n## 商业模式\n\n- 免费工具，通过开发者工具品牌赞助变现\n- 后期可做 B2B（帮企业发现技术趋势）', '[]', 'submitted', 20, NOW() - INTERVAL '15 days', NOW() - INTERVAL '16 days', NOW() - INTERVAL '15 days'),

(7, 23, E'# StarScope — 纯技术方案\n\n## 架构\n\n- **数据采集**: GitHub API v4 (GraphQL) 批量同步用户 Star\n- **Release 监控**: GitHub Webhook + Cron fallback\n- **Changelog 摘要**: LLM 自动从 Release Notes 提取关键变更\n- **通知分发**: Resend (邮件) + Slack Webhook\n- **存储**: PostgreSQL + Redis (通知队列)\n- **部署**: Vercel + Railway\n\n## API 限流应对\n\nGitHub API 限制 5000 req/hour。策略：\n- 增量同步（只拉取上次同步后的变更）\n- 热门项目共享缓存（1000 Star+ 项目全局监控）\n- 用户级按优先级调度', '[]', 'submitted', 12, NOW() - INTERVAL '14 days', NOW() - INTERVAL '15 days', NOW() - INTERVAL '14 days');

-- Idea 8: 个人碳足迹追踪 (3 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(8, 4, E'# CarbonLite — 游戏化碳足迹追踪\n\n## 1. 核心设计：让减碳像记账一样简单\n\n**每日碳预算**制：每人每天碳预算 8kg CO₂e（全球人均目标），像记账 App 一样追踪消耗。\n\n## 2. 游戏化机制\n\n| 机制 | 描述 |\n|------|------|\n| 碳预算 | 每日 8kg 预算，超支标红 |\n| 挑战任务 | "本周骑车通勤 3 天" = -5kg |\n| 好友排行 | 和朋友比谁更环保 |\n| 成就徽章 | "连续低碳 30 天""碳中和月" |\n| 种树奖励 | 累计减碳 100kg = 种一棵真树 |\n\n## 3. 数据采集方案\n\n| 类别 | 方式 | 精度 |\n|------|------|------|\n| 交通 | 接入地图 API，自动识别出行方式 | 高 |\n| 饮食 | 拍照识别食物 → 查碳排数据库 | 中 |\n| 购物 | 对接电商订单 API | 中 |\n| 能源 | 连接智能电表/手动输入 | 低 |\n\n## 4. 技术\n\n- React Native App\n- Go 后端 + PostgreSQL\n- 高德地图 API（出行检测）\n- 食物图像识别（MobileNet fine-tuned）\n- 碳排因子数据库（IPCC 2024 数据）\n\n## 5. 商业模式\n\n- Free: 基础追踪\n- Premium ¥12/月: 好友排行 + 挑战 + 详细报告\n- B2B: 企业碳足迹报告 SaaS\n- 碳中和品牌合作（赞助挑战活动）\n\n## 6. 指标\n\n| KPI | M3 | M6 |\n|-----|-----|-----|\n| MAU | 5,000 | 20,000 |\n| 日均打卡率 | 30% | 50% |\n| 平均减碳 | 5% | 15% |', '[{"step":"gamification","decision":"碳预算+种树","reason":"纯数据追踪留存低，游戏化和公益结合提升粘性"},{"step":"data","decision":"自动采集优先","reason":"手动输入坚持率<10%，必须尽量自动化"}]', 'submitted', 27, NOW() - INTERVAL '15 days', NOW() - INTERVAL '16 days', NOW() - INTERVAL '15 days'),

(8, 11, E'# GreenStep — 极简碳足迹\n\n## 核心理念\n\n碳计算太复杂，用户不需要精确到克的数据。他们需要的是**方向感**：今天是低碳日还是高碳日？这个月比上个月好还是差？\n\n## 设计：每天 3 个问题\n\n1. 今天怎么通勤的？（步行/骑车/公交/开车）\n2. 今天吃了什么？（素食/少肉/正常/大餐）\n3. 有额外消费吗？（无/小件/大件）\n\n基于回答给出评级：🟢 低碳日 / 🟡 中碳日 / 🔴 高碳日\n\n## 差异化\n\n**不追求精确，追求坚持**。3 个问题 10 秒完成 → 每日打卡率高 → 长期行为改变。\n\n## 月度报告\n\n- 低碳日占比趋势\n- 与上月对比\n- 最大改善空间（"如果每周少开一天车，月均减碳 12%"）\n\n## 合作方向\n\n- 碳中和品牌冠名挑战\n- 蚂蚁森林/种树公益联动\n- 企业 ESG 报告数据源', '[]', 'submitted', 19, NOW() - INTERVAL '14 days', NOW() - INTERVAL '15 days', NOW() - INTERVAL '14 days'),

(8, 24, E'# FootprintTracker — 微信小程序方案\n\n## 技术方案\n\n**微信小程序 MVP**（降低使用门槛，无需下载 App）\n\n### 架构\n- 前端: 微信小程序 (Taro 跨端框架)\n- 后端: Go + Echo + PostgreSQL\n- 地图: 高德地图微信 SDK\n- AI: 食物图像识别 (腾讯云 AI)\n\n### 碳排计算模型\n- 交通: 距离 × 交通方式碳因子（高德路径规划 API）\n- 饮食: 食物种类 × 碳排因子数据库（500+ 食物条目）\n- 消费: 品类平均碳排估算\n- 数据源: IPCC 排放因子 + 中国生态环境部数据\n\n### 数据准确性\n- MVP 精度目标: ±30%（足够给出方向性建议）\n- 后续优化: 机器学习模型根据用户反馈持续校准\n\n## 定价\n\n- 基础版免费\n- Pro ¥8/月: 详细报告 + 好友对比', '[]', 'submitted', 11, NOW() - INTERVAL '13 days', NOW() - INTERVAL '14 days', NOW() - INTERVAL '13 days');

-- Idea 9-13: 3-4 contributions each
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
-- Idea 9: 技术写作 AI 助手 (4 contributions)
(9, 2, E'# WriteDev — 懂代码的技术写作助手\n\n## 1. 核心差异\n\n通用 AI 写作工具（Jasper/Copy.ai）不理解代码。WriteDev **读懂代码上下文**，然后生成精准的技术文档。\n\n## 2. 核心功能\n\n| 功能 | 描述 | 场景 |\n|------|------|------|\n| 代码→文档 | 粘贴函数/模块 → 生成使用文档 | API 文档、SDK 指南 |\n| 多格式输出 | 同一内容输出为博客/README/API doc | 一次写作多处复用 |\n| 风格学习 | 上传 3 篇你的文章 → 学习你的写作风格 | 保持一致的个人品牌 |\n| SEO 优化 | 自动建议关键词、结构化数据、meta description | 技术博客 SEO |\n| 多语言 | 中英文互译，保留代码块不翻译 | 面向国际开发者 |\n\n## 3. 产品形态\n\n- **VS Code 插件**：选中代码 → 右键 → "Generate Doc"\n- **Web 编辑器**：Notion-like 体验，Markdown 实时预览\n- **CLI**：`writedev generate --input src/ --output docs/`\n\n## 4. 技术栈\n\n- VS Code Extension API + Language Server Protocol\n- GPT-4 Turbo（代码理解 + 文档生成）\n- AST 解析器（TypeScript/Python/Go）提取函数签名和类型\n- Next.js Web 编辑器\n\n## 5. 竞品\n\n| 产品 | 定位 | 代码理解 |\n|------|------|----------|\n| Mintlify Writer | 注释生成 | 基础 |\n| Swimm | 代码文档 | 中等 |\n| **WriteDev** | 完整技术写作 | 深度 |\n\n## 6. 定价\n\n- Free: 5 篇/月 + 基础功能\n- Pro $12/月: 无限 + 风格学习 + SEO + 多语言\n- Team $30/月: 团队风格规范 + 审核流程', '[{"step":"form","decision":"VS Code 插件优先","reason":"开发者在 IDE 中写代码时产生文档需求"},{"step":"tech","decision":"AST 解析+LLM","reason":"纯 LLM 对代码结构理解不够，AST 提供精确的类型和签名信息"}]', 'submitted', 36, NOW() - INTERVAL '13 days', NOW() - INTERVAL '14 days', NOW() - INTERVAL '13 days'),
(9, 7, E'# CodeScribe — 代码到文档的翻译器\n\n## 定位\n\n不是通用写作助手，是**代码到文档的专业翻译器**。输入代码，输出人能读懂的文档。\n\n## 核心功能\n\n1. **API 文档生成**：扫描路由文件 → 自动生成 OpenAPI spec + 可读文档\n2. **代码注释生成**：分析函数逻辑 → 生成 JSDoc/GoDoc/Docstring\n3. **CHANGELOG 自动化**：对比两个 Git tag → 生成人类可读的更新日志\n4. **多语言翻译**：中英文文档互转，代码块保持不变\n\n## 使用方式\n\n```bash\n# CLI\ncodescribe api --input ./routes --output ./docs/api.md\ncodescribe changelog --from v1.0 --to v2.0\n\n# VS Code\n选中函数 → Cmd+Shift+D → 生成文档\n```\n\n## 技术\n\n- Tree-sitter 解析 AST（支持 20+ 语言）\n- GPT-4 Turbo 生成自然语言\n- Git diff 分析 changelog\n\n## 定价\n\n- 开源 CLI（基础功能）\n- Cloud $8/月: 高级功能 + API + 团队协作', '[]', 'submitted', 28, NOW() - INTERVAL '12 days', NOW() - INTERVAL '13 days', NOW() - INTERVAL '12 days'),
(9, 15, E'# DocBot — CI/CD 驱动的文档自动化\n\n## 差异化\n\n**集成到 CI/CD，每次 PR 合并自动更新文档。**\n\n其他工具需要手动触发，DocBot 是全自动的：代码一变，文档跟着变。\n\n## 使用场景\n\n| 触发事件 | DocBot 动作 |\n|----------|-------------|\n| 新模块 PR 合并 | 自动生成设计文档骨架 |\n| API 路由变更 | 自动更新 API 文档 |\n| Bug 修复 PR | 生成 postmortem 模板 |\n| 版本发布 | 自动生成 Release Notes |\n\n## 技术\n\n- GitHub App（Webhook 监听 PR/Push 事件）\n- 差异分析引擎（AST diff → 识别哪些文档需要更新）\n- GPT-4 生成/更新文档\n- 自动提 PR 到 docs/ 目录\n\n## 定价\n\n- Free: 公开仓库\n- Pro $15/月: 私有仓库 + 自定义模板', '[]', 'submitted', 20, NOW() - INTERVAL '11 days', NOW() - INTERVAL '12 days', NOW() - INTERVAL '11 days'),
(9, 25, E'# TechPen — 让写技术文章像发推一样简单\n\n## 产品思路\n\n降低技术写作门槛。很多开发者有好的技术分享，但觉得"写文章太麻烦了"。\n\n## 功能\n\n1. **语音转文章**: 对着麦克风讲 5 分钟技术分享 → AI 整理成文章\n2. **代码截图美化**: 像 Carbon.sh 但集成在编辑器里\n3. **一键多平台发布**: 掘金/Medium/Dev.to/个人博客同步发布\n4. **SEO 自动优化**: 标题建议、关键词、meta\n\n## 技术\n\n- Web 编辑器 (Tiptap/ProseMirror)\n- Whisper + GPT-4 语音转文章 pipeline\n- 各平台 API 发布集成\n\n## 定价\n\n- Free: 3 篇/月\n- Pro $6/月: 无限 + 语音 + 多平台', '[]', 'submitted', 13, NOW() - INTERVAL '10 days', NOW() - INTERVAL '11 days', NOW() - INTERVAL '10 days'),

-- Idea 10: 开源项目赞助匹配 (3 contributions)
(10, 5, E'# FundOSS — 开源赞助双边市场\n\n## 1. 双边市场设计\n\n**供给侧（开源维护者）**：\n- 注册项目，展示 GitHub Star/Download/依赖数等影响力数据\n- 设置赞助目标和用途说明\n- 展示赞助回报（logo 露出、优先 issue 响应等）\n\n**需求侧（企业）**：\n- 扫描 package.json/go.mod → 自动列出使用的开源项目\n- 按技术栈和预算筛选赞助对象\n- 一键批量赞助，统一发票\n\n## 2. 匹配算法\n\n```\n企业依赖分析 → 提取所有直接依赖的开源项目\n→ 按使用频率/关键程度评分\n→ 推荐赞助金额（依赖关键度 × 项目资金缺口）\n→ 一键赞助 top N 项目\n```\n\n## 3. 商业模式\n\n- 企业端撮合服务费 10%（赞助金额的 10%）\n- 赞助效果报告（季度 PDF：赞助了哪些项目、项目的更新情况）\n- 企业 CSR 证书（"本企业赞助了 N 个开源项目"）\n\n## 4. 竞品\n\n| 产品 | 模式 | 短板 |\n|------|------|------|\n| GitHub Sponsors | 个人赞助 | 企业流程复杂，无自动推荐 |\n| Open Collective | 财务透明 | 无匹配算法，需手动寻找 |\n| Tidelift | 企业订阅 | 只覆盖部分大型项目 |\n\n## 5. MVP\n\n先做 npm 生态，扫描 package.json → 推荐 → GitHub Sponsors 支付。', '[{"step":"market","decision":"企业端切入","reason":"个人赞助金额小且不稳定，企业有预算且可规模化"},{"step":"ecosystem","decision":"npm 生态优先","reason":"JS 生态依赖数量最多，数据最丰富"}]', 'submitted', 30, NOW() - INTERVAL '11 days', NOW() - INTERVAL '12 days', NOW() - INTERVAL '11 days'),
(10, 13, E'# OpenSponsor — 让开源赞助像买 SaaS 一样简单\n\n## 核心价值\n\n企业知道应该赞助开源，但不知道赞助谁、赞助多少。OpenSponsor 把这个决策自动化。\n\n## 功能\n\n1. **依赖扫描**: 上传 package.json/Gemfile/go.mod → 看你用了哪些开源项目\n2. **赞助建议**: 按依赖使用频率推荐赞助金额\n3. **一键赞助**: 统一支付入口，自动按比例分配给各项目\n4. **赞助报告**: 年度开源贡献证书 + 影响力数据\n5. **税务支持**: 自动生成赞助发票和税务凭证\n\n## 用户旅程\n\n```\n上传 package.json → 发现你依赖 142 个开源项目\n→ 推荐每月赞助 $200 → 分配方案预览\n→ 一键支付 → 每季度收到赞助报告\n```\n\n## 定价\n\n- 企业端: 赞助金额的 8% 平台费\n- 开源端: 免费（供给侧免费是双边市场标准做法）', '[]', 'submitted', 22, NOW() - INTERVAL '10 days', NOW() - INTERVAL '11 days', NOW() - INTERVAL '10 days'),
(10, 26, E'# PatchFund — 开源悬赏平台\n\n## 创新点\n\n不做赞助，做**开源悬赏（Bounty）**。企业不是"捐钱"，而是"花钱解决问题"。\n\n## 机制\n\n1. 企业在依赖的开源项目上发布 Bounty（"实现 X 功能 → $500"）\n2. 开发者认领 → 提交 PR\n3. 项目维护者审核通过\n4. PatchFund 自动打款给开发者\n\n## 价值\n\n- **企业**: 花钱解决实际问题（不是捐赠，是采购）\n- **开发者**: 靠技能赚钱，不用"乞讨"赞助\n- **开源项目**: 获得高质量贡献\n\n## 技术\n\n- GitHub App（关联 Issue/PR）\n- Stripe Connect（分账打款）\n- 智能合约（可选，保证资金托管安全）\n\n## 商业模式\n\n- 平台佣金 15%（从 Bounty 金额中扣除）', '[]', 'submitted', 14, NOW() - INTERVAL '9 days', NOW() - INTERVAL '10 days', NOW() - INTERVAL '9 days'),

-- Idea 11: 城市探索盲盒 (3 contributions)
(11, 6, E'# CityBlind — 城市探索盲盒\n\n## 1. 核心体验\n\n每周收到一个**神秘坐标**。你不知道那里是什么——可能是一家隐藏在胡同里的咖啡馆，可能是一个没人知道的屋顶花园，也可能是一个只有本地人才去的早餐摊。\n\n**设计理念**：在确定性过剩的时代，制造惊喜。\n\n## 2. 内容来源\n\n| 来源 | 占比 | 质量保证 |\n|------|------|----------|\n| 本地博主投稿 | 40% | 实地验证 + 编辑审核 |\n| 用户打卡推荐 | 35% | 评分筛选（4.5+ 星才入库） |\n| 编辑精选 | 25% | 专业团队线下调研 |\n\n## 3. 盲盒类型\n\n- **美食盲盒**: 隐藏餐厅/小吃摊\n- **自然盲盒**: 城市里的绿洲/观景点\n- **文化盲盒**: 独立书店/艺术空间/手作工坊\n- **怀旧盲盒**: 有故事的老街/老店\n\n## 4. 商业模式\n\n- Free: 每月 1 个免费盲盒\n- Premium ¥15/月: 每周 1 个 + 专属类型选择\n- 商家合作: 新店引流（但绝不推送广告，保持"惊喜感"）\n- 品牌联名: "XX啤酒 × CityBlind 夏日探索季"\n\n## 5. MVP\n\n微信小程序，先做北京和上海，各准备 100 个地点。', '[{"step":"content","decision":"PGC+UGC 混合","reason":"纯UGC质量不可控，纯PGC不可扩展"},{"step":"city","decision":"北京+上海首发","reason":"一线城市隐藏地点密度最高"}]', 'submitted', 25, NOW() - INTERVAL '10 days', NOW() - INTERVAL '11 days', NOW() - INTERVAL '10 days'),
(11, 16, E'# Wanderly — 本地人视角的城市探索\n\n## 差异化\n\n不推荐热门景点和网红店，**专门推荐「只有当地人才知道的地方」**。\n\n## 内容策略\n\n- 每个城市招募 10-20 个"城市向导"（本地生活 5 年+）\n- 向导提交地点需包含：为什么推荐、最佳时间、隐藏技巧\n- 所有地点必须满足："如果带外地朋友来，你会带他去"\n\n## 社交元素\n\n1. **探索打卡**: 去过的地方点亮地图\n2. **区域排行榜**: 谁探索了最多角落\n3. **成就系统**: 美食猎人/公园达人/胡同专家/夜行侠\n4. **探索日记**: 分享你的盲盒体验\n\n## 增长策略\n\n打卡笔记分享到朋友圈/小红书 → 制造"这是哪里？我也想去"的好奇心 → 裂变。\n\n## 定价\n\n- Free: 每月 2 个推荐\n- 会员 ¥12/月: 每周推荐 + 独享地点', '[]', 'submitted', 18, NOW() - INTERVAL '9 days', NOW() - INTERVAL '10 days', NOW() - INTERVAL '9 days'),
(11, 27, E'# UrbanBox — 技术方案\n\n## 技术架构\n\n- **前端**: 微信小程序（Taro 跨端）\n- **后端**: Go + Echo + PostgreSQL + PostGIS\n- **地图**: 高德地图 SDK\n- **推荐算法**: LBS + 用户兴趣标签 + 协同过滤\n\n## 数据模型\n\n```sql\nspots (id, city, category, lat, lng, name, description, photos, rating)\nuser_preferences (user_id, categories[], explored_spots[])\nassignments (user_id, spot_id, week, opened_at, visited_at)\n```\n\n## 内容来源\n\n- 大众点评 POI 数据（低分筛选"小众"店铺）\n- 小红书 API 抓取低热度高评分内容\n- 用户投稿 + 人工审核\n\n## MVP: 4 周\n\n- W1: 小程序框架 + 地图基础\n- W2: 后端 API + 地点管理后台\n- W3: 盲盒分配算法 + 推送\n- W4: 打卡流程 + 测试', '[]', 'submitted', 10, NOW() - INTERVAL '8 days', NOW() - INTERVAL '9 days', NOW() - INTERVAL '8 days'),

-- Idea 12: 代码审查学习平台 (3 contributions)
(12, 3, E'# ReviewSchool — 从真实 PR 学习 Code Review\n\n## 1. 产品设计\n\n精选开源项目中的高质量 PR review 案例，打造**交互式 code review 训练平台**。\n\n## 2. 学习模式\n\n```\n用户看到一个真实 PR diff\n→ 自己写 review comments\n→ 提交后看专家/社区的实际 review\n→ AI 对比分析：你漏了什么？你多虑了什么？\n→ 获得评分和改进建议\n```\n\n## 3. 内容生产 Pipeline\n\n| 步骤 | 方法 |\n|------|------|\n| 采集 | GitHub API 爬取 10k Star+ 项目的 PR（有 3+ review comments 的） |\n| 筛选 | AI 评估 review 质量（是否有实质性讨论，非 LGTM） |\n| 分类 | 按语言/主题/难度自动标签 |\n| 标注 | 社区专家补充教学说明 |\n\n## 4. 分级系统\n\n| 级别 | 内容 | 示例 |\n|------|------|------|\n| L1 入门 | 代码风格/命名规范 | 变量命名不清晰 |\n| L2 基础 | Bug 和逻辑错误 | 边界条件未处理 |\n| L3 进阶 | 性能和安全问题 | N+1 查询/SQL 注入 |\n| L4 高级 | 设计模式和架构 | 职责不清/耦合过重 |\n| L5 专家 | 系统设计 review | 分布式一致性/可扩展性 |\n\n## 5. 商业模式\n\n- Free: L1-L2，每天 3 题\n- Pro $8/月: 全等级 + AI 分析 + 学习路径\n- Team $25/月: 团队训练 + review 水平报告', '[{"step":"content","decision":"真实 PR 而非人造案例","reason":"真实案例更有说服力，学习效果更好"}]', 'submitted', 32, NOW() - INTERVAL '9 days', NOW() - INTERVAL '10 days', NOW() - INTERVAL '9 days'),
(12, 10, E'# CodeReviewGym — 代码审查刻意练习\n\n## 练习模式\n\n**三种训练模式**：\n\n1. **找 Bug 模式**: 代码中故意藏了 1-3 个问题，找出来\n   - "这段代码有一个并发安全问题，找到它"\n2. **改进模式**: 代码可以运行，但写得不好，提优化建议\n   - "这个函数有 50 行，如何重构？"\n3. **设计评估**: 给出两个实现方案，评估哪个更好\n   - "方案 A 用继承，方案 B 用组合，哪个更合适？"\n\n## 分级\n\nL1 语法 → L2 逻辑 → L3 性能 → L4 设计 → L5 架构\n\n每个级别 50+ 题，覆盖 Python/JavaScript/Go/Java。\n\n## 游戏化\n\n- 段位系统（铜→银→金→钻石）\n- 每日挑战 + 排行榜\n- 连续打卡奖励\n\n## 技术\n\n- Next.js + Monaco Editor (diff viewer)\n- 题库存储在 PostgreSQL\n- AI 评分引擎（对比用户 review 与标准答案）', '[]', 'submitted', 23, NOW() - INTERVAL '8 days', NOW() - INTERVAL '9 days', NOW() - INTERVAL '8 days'),
(12, 20, E'# PRMaster — AI 驱动的 Review 教学\n\n## 思路\n\n用真实的开源 PR 作为教材，AI 老师逐行讲解每条 review comment 背后的逻辑。\n\n## 学习流程\n\n1. 打开一个真实的 GitHub PR diff\n2. 用户尝试找出问题并写 review\n3. 点击"显示专家 review"\n4. AI 逐条解释：\n   - 为什么这里需要 review？\n   - review 的依据是什么原则？\n   - 类似问题在其他项目中的处理方式\n\n## 技术\n\n- GitHub API 抓取 PR diff + comments\n- GPT-4 分析 review 逻辑并生成教学说明\n- 交互式 diff viewer（高亮 review 位置）\n- 向量搜索（找相似的 review 案例）\n\n## 定价\n\n- Free: 每天 2 个 PR\n- Pro $10/月: 无限 + AI 教学 + 个人弱点分析', '[]', 'submitted', 15, NOW() - INTERVAL '7 days', NOW() - INTERVAL '8 days', NOW() - INTERVAL '7 days'),

-- Idea 13: API 变更通知服务 (3 contributions)
(13, 4, E'# APIWatch — 第三方 API 变更监控与预警\n\n## 1. 监控四层机制\n\n| 层级 | 方法 | 覆盖范围 |\n|------|------|----------|\n| 文档 Diff | 定时抓取 API 文档页面 + diff 对比 | 有文档的 API |\n| OpenAPI Spec | 解析 spec 文件变更（endpoint/参数/schema） | 提供 spec 的 API |\n| GitHub Release | 监控 SDK 仓库的 Release 和 CHANGELOG | 开源 SDK |\n| 行为监测 | 定时发测试请求，对比响应结构 | 所有 API |\n\n## 2. 通知分级\n\n| 级别 | 类型 | 通知方式 |\n|------|------|----------|\n| P0 | Breaking change | 即时 Slack/邮件告警 |\n| P1 | Deprecation notice | 周报 |\n| P2 | 新功能/新 endpoint | 月报 |\n| Info | 文档小修改 | Dashboard 记录 |\n\n## 3. 迁移辅助\n\n检测到 breaking change 后：\n1. AI 分析影响范围（你的代码哪些地方调用了变更的 API）\n2. 自动生成迁移指南\n3. 创建 GitHub Issue 追踪修复进度\n\n## 4. 技术\n\n- Go 爬虫引擎（定时任务 + Webhook）\n- Diff 引擎（文本/JSON/YAML 多格式对比）\n- OpenAI API（changelog 摘要 + 迁移指南生成）\n- PostgreSQL + Redis\n\n## 5. 定价\n\n- Free: 监控 3 个 API\n- Pro $15/月: 无限 API + AI 迁移辅助 + Slack 集成\n- Enterprise $50/月: CI/CD 集成 + 团队 Dashboard', '[{"step":"coverage","decision":"四层监控","reason":"单一方式覆盖不全，文档可能不更新但API已变更"},{"step":"priority","decision":"P0 即时告警","reason":"breaking change 如果不及时发现会导致线上故障"}]', 'submitted', 29, NOW() - INTERVAL '8 days', NOW() - INTERVAL '9 days', NOW() - INTERVAL '8 days'),
(13, 12, E'# ChangeGuard — CI/CD 集成的 API 变更检测\n\n## 核心思路\n\n**不做独立产品，做 CI/CD 插件。** 每次构建时自动检查你依赖的 API 是否有变更。\n\n## GitHub Action 集成\n\n```yaml\n- uses: changeguard/check@v1\n  with:\n    apis: stripe,twilio,github\n    severity: high\n```\n\n## 功能\n\n1. 每次 CI 运行时检查依赖 API 变更\n2. 变更影响评估（High/Medium/Low）\n3. 发现 breaking change → CI 警告（不阻塞，但显著提醒）\n4. 自动创建 Issue 追踪修复进度\n5. 支持 REST API + GraphQL + gRPC\n\n## 技术\n\n- GitHub Action / GitLab CI 插件\n- 中央 API 变更数据库（由 ChangeGuard 维护）\n- Webhook 通知\n\n## 定价\n\n- Free: 公开仓库\n- Pro $10/月: 私有仓库 + 高级通知', '[]', 'submitted', 21, NOW() - INTERVAL '7 days', NOW() - INTERVAL '8 days', NOW() - INTERVAL '7 days'),
(13, 28, E'# APIDrift — 基于行为的 API 变更检测\n\n## 独特方案\n\n不监控文档，监控**实际 API 行为**。\n\n### 原理\n- 注册你使用的 API endpoint\n- APIDrift 定期发送测试请求\n- 对比响应的 JSON 结构（字段增减、类型变化、状态码变化）\n- 发现差异 → 报警\n\n### 优势\n- 发现文档没写但实际已变更的"暗变更"\n- 比爬文档更可靠（有些 API 根本不更新文档）\n- 可以检测到性能退化（响应时间变慢）\n\n### 局限\n- 需要有效的测试 credentials\n- 有些 API 会限制测试请求频率\n- 无法检测"新增 endpoint"（因为不知道新的 URL）\n\n## 技术\n\n- Go 调度器 + 分布式测试节点\n- JSON Schema diff 引擎\n- Grafana 集成（响应时间监控）\n\n## 定价\n\n- Free: 5 endpoint\n- Pro $12/月: 无限 + 高级对比 + 历史记录', '[]', 'submitted', 11, NOW() - INTERVAL '6 days', NOW() - INTERVAL '7 days', NOW() - INTERVAL '6 days');

-- Open ideas contributions (mix of draft and submitted)
-- Idea 14: AI Commit Message (4 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(14, 1, E'# CommitCraft — 智能 Commit Message 生成器\n\n## 方案概要\n\n分析 `git diff` 的语义内容，自动生成符合 Conventional Commits 规范的 commit message。\n\n## 核心功能\n\n| 功能 | 描述 |\n|------|------|\n| Diff 语义分析 | 理解代码变更的意图（新功能/修复/重构） |\n| Conventional Commits | 自动选择 feat/fix/refactor/docs 前缀 |\n| 中英文支持 | 根据项目设置生成对应语言 |\n| Scope 推断 | 根据文件路径推断变更范围 |\n| Breaking Change 检测 | 自动标记 BREAKING CHANGE |\n\n## 产品形态\n\n- **CLI**: `commitcraft` → 分析 staged changes → 生成 message → 确认提交\n- **VS Code 插件**: Git 面板集成，一键生成\n- **Git Hook**: `prepare-commit-msg` hook 自动触发\n\n## 技术\n\n- `git diff --staged` 获取变更\n- GPT-4 Turbo 分析 diff 语义\n- 本地缓存项目历史 commit 风格\n- Token 优化: 大 diff 自动截断到关键变更\n\n## 定价\n\n- 开源 CLI（本地 LLM 支持）\n- Cloud $5/月: GPT-4 驱动 + 团队风格同步', '[]', 'submitted', 8, NOW() - INTERVAL '1 day', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day'),
(14, 5, E'# GitMsg AI — 学习项目风格的 Commit 助手\n\n## 差异化\n\n不只是生成 commit message，而是**学习项目的历史 commit 风格**后生成风格一致的 message。\n\n## 工作原理\n\n```\n1. 首次使用: 分析最近 200 条 commit message\n   → 提取风格模式（前缀/语言/长度/大小写）\n2. 每次 commit: 分析 diff + 匹配项目风格\n   → 生成风格一致的 message\n```\n\n## 示例\n\n如果项目风格是 `[模块] 动词 描述`：\n- `[auth] 修复 token 过期判断逻辑`\n- `[api] 添加用户资料接口`\n\n如果项目风格是 Conventional Commits：\n- `fix(auth): correct token expiry check`\n- `feat(api): add user profile endpoint`\n\n## 技术\n\n- Few-shot learning: 用项目历史 commits 作为 examples\n- 本地 embedding: 风格向量化存储\n- 支持 Ollama 本地模型（隐私友好）', '[]', 'submitted', 5, NOW() - INTERVAL '12 hours', NOW() - INTERVAL '1 day', NOW() - INTERVAL '12 hours'),
(14, 18, E'# Commit AI 草稿\n\n## 初步想法\n\n想做一个 commit message 生成器，但还在调研技术方案...\n\n### 待研究\n- [ ] 本地 LLM vs API\n- [ ] diff 太大怎么处理\n- [ ] 多语言支持策略\n- [ ] 如何处理 monorepo', '[]', 'draft', 0, NULL, NOW() - INTERVAL '1 day', NOW() - INTERVAL '6 hours'),
(14, 22, E'# AutoCommit — Git Hook 自动化方案\n\n## 核心功能\n\n- **Git Hook 集成**: `prepare-commit-msg` hook，commit 时自动填充 message\n- **多风格模板**: Conventional Commits / Angular / 自定义模板\n- **团队配置**: `.autocommitrc` 文件统一团队规范\n- **Review 模式**: 生成后在编辑器中打开，允许修改\n\n## 安装\n\n```bash\nnpm install -g autocommit\nautocommit init  # 安装 git hook + 生成配置文件\n```\n\n## 技术\n\n- Node.js CLI\n- OpenAI / Ollama API\n- `.autocommitrc` JSON 配置', '[]', 'submitted', 3, NOW() - INTERVAL '6 hours', NOW() - INTERVAL '1 day', NOW() - INTERVAL '6 hours');

-- Idea 15: 开发者人体工学 (3 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(15, 7, E'# DevHealth — 开发者健康提醒工具\n\n## 方案\n\n**番茄钟 + 人体工学提醒的结合体**，专为程序员设计。\n\n## 核心功能\n\n| 定时 | 提醒内容 |\n|------|----------|\n| 每 25 分钟 | 番茄钟结束，休息 5 分钟 |\n| 每 45 分钟 | 站起来做 30 秒拉伸（附动画指导） |\n| 每 2 小时 | 20-20-20 护眼（看 20 英尺外 20 秒） |\n| 每 4 小时 | 建议喝水 + 5 分钟散步 |\n\n## 产品形态\n\n- macOS 菜单栏常驻（极简 UI）\n- 提醒方式: 系统通知 + 可选全屏遮罩（强制休息模式）\n- 统计: 每日/周编码时长、休息次数、健康评分\n\n## 技术\n\n- Swift (macOS native)\n- 后续: Electron 跨平台\n- 监控键盘/鼠标活动检测"是否在编码"\n\n## 差异化\n\n市面番茄钟很多，但没有**理解程序员工作节奏**的：\n- 检测到你在 debug → 不打断（延后提醒）\n- 检测到你刚 push 完 → 适合休息的好时机\n\n## 定价\n\n- Free: 基础定时提醒\n- Pro $3/月: 智能检测 + 统计 + 团队健康报告', '[]', 'submitted', 6, NOW() - INTERVAL '2 days', NOW() - INTERVAL '3 days', NOW() - INTERVAL '2 days'),
(15, 11, E'# PostureGuard — 无摄像头姿态检测\n\n## 创新点\n\n**不用摄像头**，通过键盘打字节奏和鼠标移动模式推断疲劳程度和姿态问题。\n\n## 原理\n\n| 信号 | 含义 |\n|------|------|\n| 打字速度下降 30%+ | 可能疲劳 |\n| 打字错误率上升 | 注意力下降 |\n| 鼠标移动变慢/不精准 | 手臂疲劳 |\n| 连续编码 >90 分钟无停顿 | 需要强制提醒 |\n\n## 隐私优势\n\n- 不需要摄像头权限\n- 不记录击键内容（只分析节奏）\n- 所有数据本地处理\n\n## 技术\n\n- macOS Accessibility API（键盘/鼠标事件监控）\n- 本地 ML 模型（疲劳检测）\n- Swift native app\n\n## 风险\n\n- 准确率可能不够高 → 需要大量数据训练\n- 用户可能觉得被监控 → 强调隐私设计', '[]', 'submitted', 4, NOW() - INTERVAL '1 day', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day'),
(15, 29, E'# 健康码 — 草稿\n\n## 初步思路\n\n从 Apple Watch 健康数据切入...\n\n### 想法\n- 读取 Apple Watch 心率、站立时间、活动量\n- 结合 IDE 使用时间（从 WakaTime API 获取）\n- 生成"开发者健康评分"\n\n### 待确认\n- [ ] Apple HealthKit API 权限限制\n- [ ] WakaTime API 是否可用\n- [ ] 隐私合规性', '[]', 'draft', 0, NULL, NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day');

-- Idea 16: 技术播客推荐 (2 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(16, 9, E'# PodMatch — 技术播客精准推荐引擎\n\n## 推荐算法（三层架构）\n\n### Layer 1: 标签匹配\n用户填写技术栈标签（React/Go/DevOps/AI...）→ 匹配播客节目标签\n\n### Layer 2: 协同过滤\n"和你技术栈相似的开发者还在听什么"\n- User-based CF: 技术栈重合度 > 70% 的用户\n- Item-based CF: 听了 A 播客的人也听了 B\n\n### Layer 3: 内容分析\n- NLP 提取每集主题关键词\n- 匹配用户近期关注的技术话题\n- 例: 你最近在写 Rust → 推荐 Rust 相关单集\n\n## 数据源\n\n- Apple Podcasts API\n- 小宇宙开放接口\n- RSS Feed 爬取（Podcast Index）\n- 播客转录 + NLP 主题提取\n\n## MVP\n\n先做 Web App，用户填技术栈 → 返回 Top 10 推荐。不做 App，不做播放器。\n\n## 定价\n\n- Free: 基础推荐\n- Pro $4/月: 个性化单集推荐 + 每周精选邮件', '[]', 'submitted', 4, NOW() - INTERVAL '12 hours', NOW() - INTERVAL '1 day', NOW() - INTERVAL '12 hours'),
(16, 24, E'# TechCast 草稿\n\n## 初步想法\n\n先做一个**人工精选的技术播客列表**，不急着做算法推荐。\n\n### 计划\n- 分类: 前端/后端/DevOps/AI/创业\n- 每个分类精选 10 个播客\n- 每周人工推荐 3 个"本周必听"\n- 后续再加 AI 推荐\n\n### 待定\n- [ ] 用什么形态？Newsletter？Web？小程序？\n- [ ] 如何获取播客元数据', '[]', 'draft', 0, NULL, NOW() - INTERVAL '18 hours', NOW() - INTERVAL '6 hours');

-- Idea 17: 开源许可证检查器 (3 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(17, 3, E'# LicenseGuard — 开源许可证合规 CI 工具\n\n## 核心功能\n\n### CI/CD 一键接入\n```yaml\n# .github/workflows/license-check.yml\n- uses: licenseguard/check@v1\n  with:\n    policy: strict  # strict|moderate|permissive\n    fail-on: copyleft  # 发现 copyleft 许可证则 CI 失败\n```\n\n### 支持的包管理器\n- npm (package-lock.json)\n- pip (requirements.txt / Pipfile)\n- Go (go.sum)\n- Cargo (Cargo.lock)\n- Maven (pom.xml)\n- Gradle (build.gradle)\n\n### 功能详情\n| 功能 | 描述 |\n|------|------|\n| 许可证检测 | 扫描所有直接+间接依赖的许可证 |\n| 冲突检测 | 检测不兼容许可证组合（如 MIT + GPL） |\n| 合规报告 | 生成 PDF 合规报告（满足审计需求） |\n| SBOM 生成 | 输出 CycloneDX / SPDX 格式 |\n| 策略配置 | 自定义允许/禁止的许可证列表 |\n\n## 技术\n\n- Go CLI 工具（快速、单二进制）\n- 许可证数据库（SPDX 标准 + 自定义规则）\n- GitHub Action / GitLab CI 插件\n\n## 定价\n\n- 开源 CLI: 免费\n- Cloud Dashboard $10/月: 历史记录 + 团队管理 + 告警', '[]', 'submitted', 7, NOW() - INTERVAL '3 days', NOW() - INTERVAL '4 days', NOW() - INTERVAL '3 days'),
(17, 10, E'# CompliBot — 合规检测 + 修复建议\n\n## 差异化\n\n不只检测问题，还**提供修复建议和替代依赖推荐**。\n\n## 工作流程\n\n```\n扫描依赖 → 发现 GPL 库\n→ 评估影响（是否传染？是否运行时依赖？）\n→ 推荐替代方案（MIT 许可的同功能库）\n→ 生成 PR 替换依赖\n```\n\n## 智能功能\n\n1. **传染性分析**: 区分"直接依赖"和"开发依赖"，dev 依赖的 GPL 通常不传染\n2. **替代推荐**: 维护一个"许可证友好的替代库"数据库\n3. **风险评分**: 综合许可证类型 × 依赖深度 × 使用方式\n\n## 技术\n\n- Node.js CLI + VS Code 插件\n- 替代库数据库（人工维护 + 社区贡献）\n- GitHub App（自动 PR）', '[]', 'submitted', 5, NOW() - INTERVAL '2 days', NOW() - INTERVAL '3 days', NOW() - INTERVAL '2 days'),
(17, 30, E'# 许可证雷达 — 依赖树可视化\n\n## 核心功能\n\n1. **依赖树可视化**: 交互式树状图，颜色标记许可证类型\n   - 绿色: MIT/Apache (安全)\n   - 黄色: LGPL (需注意)\n   - 红色: GPL/AGPL (传染性)\n2. **传染性分析**: 高亮"传染路径"（哪条依赖链引入了 copyleft）\n3. **SBOM 生成**: CycloneDX / SPDX 标准输出\n4. **变更追踪**: 新 PR 引入了新许可证 → 自动 comment 提醒\n\n## 技术\n\n- Web UI: Next.js + D3.js 可视化\n- 后端: Go 解析依赖树\n- 数据: SPDX license-list-data\n\n## 定价\n\n- Free: 公开仓库\n- Pro $8/月: 私有仓库 + 变更追踪', '[]', 'submitted', 3, NOW() - INTERVAL '1 day', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day');

-- Idea 18: 代码片段搜索 (2 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(18, 6, E'# SnippetSearch — 跨项目代码语义搜索\n\n## 问题\n\n"我记得以前写过一个 Redis 分布式锁的实现，但不记得在哪个项目里了。"\n\n## 方案\n\n### 本地索引引擎\n1. 扫描指定目录下的所有项目\n2. 对每个代码文件生成 embedding（OpenAI / 本地模型）\n3. 存入本地向量数据库\n4. 自然语言搜索\n\n### 搜索示例\n```\n$ snippet search "redis distributed lock"\n→ ~/projects/payment-service/pkg/lock/redis.go:15\n→ ~/projects/old-api/utils/cache.go:42\n```\n\n### 功能\n- 自然语言搜索（"那个处理 JWT 的中间件"）\n- 正则搜索（fallback）\n- 文件类型过滤\n- 最近使用优先排序\n\n## 技术\n\n- Rust CLI（高性能文件扫描）\n- SQLite + sqlite-vss（本地向量搜索）\n- OpenAI Embedding API / 本地 ONNX 模型\n- 增量索引（文件变更时更新）\n\n## 定价\n\n- 开源 CLI（本地模型）\n- Pro $5/月: 云端 embedding + 跨设备同步', '[]', 'submitted', 5, NOW() - INTERVAL '1 day', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day'),
(18, 15, E'# CodeVault 草稿\n\n## 初步构思\n\n做一个个人代码片段收藏夹...\n\n### 想法\n- 从 GitHub Gist 导入\n- VS Code 插件选中代码 → 保存到 vault\n- 标签分类 + 全文搜索\n\n### 还没想清楚的\n- [ ] 和 SnippetSearch 的区别是什么？\n- [ ] 手动收藏 vs 自动索引\n- [ ] 是否需要语义搜索', '[]', 'draft', 0, NULL, NOW() - INTERVAL '1 day', NOW() - INTERVAL '6 hours');


-- ============================================================
-- 4. Votes (~100, only for closed ideas with submitted contributions)
-- ============================================================

-- Idea 1 votes (contribs 1-5, 12 votes)
INSERT INTO votes (idea_id, voter_id, contribution_id, voted_at) VALUES
(1, 3, 1, NOW() - INTERVAL '29 days'), (1, 5, 1, NOW() - INTERVAL '28 days'),
(1, 6, 2, NOW() - INTERVAL '28 days'), (1, 8, 2, NOW() - INTERVAL '27 days'),
(1, 9, 2, NOW() - INTERVAL '27 days'), (1, 11, 1, NOW() - INTERVAL '26 days'),
(1, 12, 2, NOW() - INTERVAL '26 days'), (1, 13, 3, NOW() - INTERVAL '25 days'),
(1, 15, 1, NOW() - INTERVAL '25 days'), (1, 16, 2, NOW() - INTERVAL '24 days'),
(1, 17, 3, NOW() - INTERVAL '24 days'), (1, 18, 1, NOW() - INTERVAL '23 days');

-- Idea 2 votes (contribs 6-9, 10 votes)
INSERT INTO votes (idea_id, voter_id, contribution_id, voted_at) VALUES
(2, 2, 6, NOW() - INTERVAL '25 days'), (2, 3, 6, NOW() - INTERVAL '25 days'),
(2, 4, 6, NOW() - INTERVAL '24 days'), (2, 5, 7, NOW() - INTERVAL '24 days'),
(2, 7, 6, NOW() - INTERVAL '23 days'), (2, 8, 8, NOW() - INTERVAL '23 days'),
(2, 10, 6, NOW() - INTERVAL '22 days'), (2, 11, 7, NOW() - INTERVAL '22 days'),
(2, 12, 9, NOW() - INTERVAL '21 days'), (2, 14, 8, NOW() - INTERVAL '21 days');

-- Idea 3 votes (contribs 10-14, 11 votes)
INSERT INTO votes (idea_id, voter_id, contribution_id, voted_at) VALUES
(3, 1, 10, NOW() - INTERVAL '23 days'), (3, 2, 10, NOW() - INTERVAL '22 days'),
(3, 4, 10, NOW() - INTERVAL '22 days'), (3, 5, 11, NOW() - INTERVAL '21 days'),
(3, 6, 10, NOW() - INTERVAL '21 days'), (3, 7, 11, NOW() - INTERVAL '20 days'),
(3, 9, 12, NOW() - INTERVAL '20 days'), (3, 10, 10, NOW() - INTERVAL '19 days'),
(3, 12, 11, NOW() - INTERVAL '19 days'), (3, 14, 13, NOW() - INTERVAL '18 days'),
(3, 15, 10, NOW() - INTERVAL '18 days');

-- Idea 4 votes (contribs 15-18, 8 votes)
INSERT INTO votes (idea_id, voter_id, contribution_id, voted_at) VALUES
(4, 1, 15, NOW() - INTERVAL '21 days'), (4, 3, 15, NOW() - INTERVAL '20 days'),
(4, 6, 16, NOW() - INTERVAL '20 days'), (4, 7, 15, NOW() - INTERVAL '19 days'),
(4, 8, 16, NOW() - INTERVAL '19 days'), (4, 9, 15, NOW() - INTERVAL '18 days'),
(4, 10, 17, NOW() - INTERVAL '18 days'), (4, 11, 16, NOW() - INTERVAL '17 days');

-- Idea 5 votes (contribs 19-22, 9 votes)
INSERT INTO votes (idea_id, voter_id, contribution_id, voted_at) VALUES
(5, 2, 19, NOW() - INTERVAL '19 days'), (5, 3, 19, NOW() - INTERVAL '18 days'),
(5, 4, 20, NOW() - INTERVAL '18 days'), (5, 5, 19, NOW() - INTERVAL '17 days'),
(5, 6, 19, NOW() - INTERVAL '17 days'), (5, 7, 20, NOW() - INTERVAL '16 days'),
(5, 9, 21, NOW() - INTERVAL '16 days'), (5, 10, 19, NOW() - INTERVAL '15 days'),
(5, 11, 20, NOW() - INTERVAL '15 days');

-- Idea 6 votes (contribs 23-25, 7 votes)
INSERT INTO votes (idea_id, voter_id, contribution_id, voted_at) VALUES
(6, 1, 23, NOW() - INTERVAL '17 days'), (6, 2, 23, NOW() - INTERVAL '16 days'),
(6, 5, 24, NOW() - INTERVAL '16 days'), (6, 6, 23, NOW() - INTERVAL '15 days'),
(6, 7, 24, NOW() - INTERVAL '15 days'), (6, 8, 23, NOW() - INTERVAL '14 days'),
(6, 10, 25, NOW() - INTERVAL '14 days');

-- Idea 7 votes (contribs 26-29, 8 votes)
INSERT INTO votes (idea_id, voter_id, contribution_id, voted_at) VALUES
(7, 1, 26, NOW() - INTERVAL '15 days'), (7, 3, 26, NOW() - INTERVAL '14 days'),
(7, 4, 27, NOW() - INTERVAL '14 days'), (7, 5, 26, NOW() - INTERVAL '13 days'),
(7, 9, 26, NOW() - INTERVAL '13 days'), (7, 10, 27, NOW() - INTERVAL '12 days'),
(7, 11, 28, NOW() - INTERVAL '12 days'), (7, 12, 26, NOW() - INTERVAL '11 days');

-- Idea 8 votes (contribs 30-32, 6 votes)
INSERT INTO votes (idea_id, voter_id, contribution_id, voted_at) VALUES
(8, 1, 30, NOW() - INTERVAL '13 days'), (8, 2, 30, NOW() - INTERVAL '12 days'),
(8, 3, 31, NOW() - INTERVAL '12 days'), (8, 5, 30, NOW() - INTERVAL '11 days'),
(8, 6, 31, NOW() - INTERVAL '11 days'), (8, 7, 30, NOW() - INTERVAL '10 days');

-- Idea 9 votes (contribs 33-36, 7 votes)
INSERT INTO votes (idea_id, voter_id, contribution_id, voted_at) VALUES
(9, 1, 33, NOW() - INTERVAL '11 days'), (9, 3, 33, NOW() - INTERVAL '10 days'),
(9, 4, 34, NOW() - INTERVAL '10 days'), (9, 5, 33, NOW() - INTERVAL '9 days'),
(9, 6, 34, NOW() - INTERVAL '9 days'), (9, 8, 35, NOW() - INTERVAL '8 days'),
(9, 9, 33, NOW() - INTERVAL '8 days');

-- Idea 10 votes (contribs 37-39, 5 votes)
INSERT INTO votes (idea_id, voter_id, contribution_id, voted_at) VALUES
(10, 1, 37, NOW() - INTERVAL '9 days'), (10, 2, 37, NOW() - INTERVAL '8 days'),
(10, 3, 38, NOW() - INTERVAL '8 days'), (10, 6, 37, NOW() - INTERVAL '7 days'),
(10, 7, 38, NOW() - INTERVAL '7 days');

-- Idea 11 votes (contribs 40-42, 5 votes)
INSERT INTO votes (idea_id, voter_id, contribution_id, voted_at) VALUES
(11, 1, 40, NOW() - INTERVAL '8 days'), (11, 2, 40, NOW() - INTERVAL '7 days'),
(11, 3, 41, NOW() - INTERVAL '7 days'), (11, 4, 40, NOW() - INTERVAL '6 days'),
(11, 5, 41, NOW() - INTERVAL '6 days');

-- Idea 12 votes (contribs 43-45, 6 votes)
INSERT INTO votes (idea_id, voter_id, contribution_id, voted_at) VALUES
(12, 1, 43, NOW() - INTERVAL '7 days'), (12, 2, 43, NOW() - INTERVAL '6 days'),
(12, 4, 44, NOW() - INTERVAL '6 days'), (12, 5, 43, NOW() - INTERVAL '5 days'),
(12, 6, 44, NOW() - INTERVAL '5 days'), (12, 7, 43, NOW() - INTERVAL '4 days');

-- Idea 13 votes (contribs 46-48, 6 votes)
INSERT INTO votes (idea_id, voter_id, contribution_id, voted_at) VALUES
(13, 1, 46, NOW() - INTERVAL '6 days'), (13, 2, 46, NOW() - INTERVAL '5 days'),
(13, 3, 47, NOW() - INTERVAL '5 days'), (13, 5, 46, NOW() - INTERVAL '4 days'),
(13, 6, 47, NOW() - INTERVAL '4 days'), (13, 7, 46, NOW() - INTERVAL '3 days');


-- ============================================================
-- 5. Reveal Snapshots (for all 13 closed ideas)
-- ============================================================

-- Idea 1: contrib 1=5票, 2=4票, 3=2票, 4=0票, 5=1票 → total=12
INSERT INTO reveal_snapshots (idea_id, ranked_results, total_votes, revealed_at) VALUES
(1, '[{"contribution_id":1,"vote_count":5,"rank":1,"is_featured":true},{"contribution_id":2,"vote_count":4,"rank":2,"is_featured":true},{"contribution_id":3,"vote_count":2,"rank":3,"is_featured":true},{"contribution_id":5,"vote_count":1,"rank":4,"is_featured":false},{"contribution_id":4,"vote_count":0,"rank":5,"is_featured":false}]', 12, NOW() - INTERVAL '18 days'),

-- Idea 2: contrib 6=5票, 7=2票, 8=2票, 9=1票 → total=10
(2, '[{"contribution_id":6,"vote_count":5,"rank":1,"is_featured":true},{"contribution_id":7,"vote_count":2,"rank":2,"is_featured":true},{"contribution_id":8,"vote_count":2,"rank":2,"is_featured":true},{"contribution_id":9,"vote_count":1,"rank":4,"is_featured":false}]', 10, NOW() - INTERVAL '15 days'),

-- Idea 3: contrib 10=6票, 11=3票, 12=1票, 13=1票, 14=0票 → total=11
(3, '[{"contribution_id":10,"vote_count":6,"rank":1,"is_featured":true},{"contribution_id":11,"vote_count":3,"rank":2,"is_featured":true},{"contribution_id":12,"vote_count":1,"rank":3,"is_featured":true},{"contribution_id":13,"vote_count":1,"rank":3,"is_featured":true},{"contribution_id":14,"vote_count":0,"rank":5,"is_featured":false}]', 11, NOW() - INTERVAL '13 days'),

-- Idea 4: contrib 15=4票, 16=3票, 17=1票, 18=0票 → total=8
(4, '[{"contribution_id":15,"vote_count":4,"rank":1,"is_featured":true},{"contribution_id":16,"vote_count":3,"rank":2,"is_featured":true},{"contribution_id":17,"vote_count":1,"rank":3,"is_featured":true},{"contribution_id":18,"vote_count":0,"rank":4,"is_featured":false}]', 8, NOW() - INTERVAL '11 days'),

-- Idea 5: contrib 19=5票, 20=3票, 21=1票, 22=0票 → total=9
(5, '[{"contribution_id":19,"vote_count":5,"rank":1,"is_featured":true},{"contribution_id":20,"vote_count":3,"rank":2,"is_featured":true},{"contribution_id":21,"vote_count":1,"rank":3,"is_featured":true},{"contribution_id":22,"vote_count":0,"rank":4,"is_featured":false}]', 9, NOW() - INTERVAL '9 days'),

-- Idea 6: contrib 23=4票, 24=2票, 25=1票 → total=7
(6, '[{"contribution_id":23,"vote_count":4,"rank":1,"is_featured":true},{"contribution_id":24,"vote_count":2,"rank":2,"is_featured":true},{"contribution_id":25,"vote_count":1,"rank":3,"is_featured":true}]', 7, NOW() - INTERVAL '7 days'),

-- Idea 7: contrib 26=5票, 27=2票, 28=1票, 29=0票 → total=8
(7, '[{"contribution_id":26,"vote_count":5,"rank":1,"is_featured":true},{"contribution_id":27,"vote_count":2,"rank":2,"is_featured":true},{"contribution_id":28,"vote_count":1,"rank":3,"is_featured":true},{"contribution_id":29,"vote_count":0,"rank":4,"is_featured":false}]', 8, NOW() - INTERVAL '5 days'),

-- Idea 8: contrib 30=4票, 31=2票, 32=0票 → total=6
(8, '[{"contribution_id":30,"vote_count":4,"rank":1,"is_featured":true},{"contribution_id":31,"vote_count":2,"rank":2,"is_featured":true},{"contribution_id":32,"vote_count":0,"rank":3,"is_featured":false}]', 6, NOW() - INTERVAL '3 days'),

-- Idea 9: contrib 33=4票, 34=2票, 35=1票, 36=0票 → total=7
(9, '[{"contribution_id":33,"vote_count":4,"rank":1,"is_featured":true},{"contribution_id":34,"vote_count":2,"rank":2,"is_featured":true},{"contribution_id":35,"vote_count":1,"rank":3,"is_featured":true},{"contribution_id":36,"vote_count":0,"rank":4,"is_featured":false}]', 7, NOW() - INTERVAL '1 day'),

-- Idea 10: contrib 37=3票, 38=2票, 39=0票 → total=5
(10, '[{"contribution_id":37,"vote_count":3,"rank":1,"is_featured":true},{"contribution_id":38,"vote_count":2,"rank":2,"is_featured":true},{"contribution_id":39,"vote_count":0,"rank":3,"is_featured":false}]', 5, NOW() - INTERVAL '12 hours'),

-- Idea 11: contrib 40=3票, 41=2票, 42=0票 → total=5
(11, '[{"contribution_id":40,"vote_count":3,"rank":1,"is_featured":true},{"contribution_id":41,"vote_count":2,"rank":2,"is_featured":true},{"contribution_id":42,"vote_count":0,"rank":3,"is_featured":false}]', 5, NOW() - INTERVAL '6 hours'),

-- Idea 12: contrib 43=4票, 44=2票, 45=0票 → total=6
(12, '[{"contribution_id":43,"vote_count":4,"rank":1,"is_featured":true},{"contribution_id":44,"vote_count":2,"rank":2,"is_featured":true},{"contribution_id":45,"vote_count":0,"rank":3,"is_featured":false}]', 6, NOW() - INTERVAL '2 hours'),

-- Idea 13: contrib 46=4票, 47=2票, 48=0票 → total=6
(13, '[{"contribution_id":46,"vote_count":4,"rank":1,"is_featured":true},{"contribution_id":47,"vote_count":2,"rank":2,"is_featured":true},{"contribution_id":48,"vote_count":0,"rank":3,"is_featured":false}]', 6, NOW() - INTERVAL '1 hour');

COMMIT;

-- Verify counts
SELECT 'users' AS entity, COUNT(*) AS count FROM users
UNION ALL SELECT 'ideas', COUNT(*) FROM ideas
UNION ALL SELECT 'contributions', COUNT(*) FROM contributions
UNION ALL SELECT 'submitted_contributions', COUNT(*) FROM contributions WHERE status = 'submitted'
UNION ALL SELECT 'votes', COUNT(*) FROM votes
UNION ALL SELECT 'reveal_snapshots', COUNT(*) FROM reveal_snapshots;
