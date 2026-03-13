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
(1, 2, E'# 宠物社交 App 产品方案\n\n## 核心价值\n让每只宠物都有自己的社交主页，主人通过宠物视角社交。\n\n## 竞品分析\n| 竞品 | 优势 | 不足 |\n|------|------|------|\n| 小红书宠物话题 | 流量大 | 非垂直场景 |\n| 波奇网 | 电商成熟 | 社交薄弱 |\n\n## 核心功能\n1. **宠物档案**：品种、年龄、性格标签\n2. **附近宠友**：基于 LBS 的宠物配对\n3. **遛弯打卡**：GPS 轨迹 + 照片分享\n4. **宠物服务**：美容、寄养预约\n\n## 用户画像\n- 95后/00后城市养宠人群\n- 月收入 8k-20k\n- 日均使用社交 App 2.5 小时\n\n## 变现模式\n- 宠物服务佣金（15%）\n- 品牌合作推广\n- 会员订阅（高级滤镜/无限配对）', '[]', 'submitted', 42, NOW() - INTERVAL '30 days', NOW() - INTERVAL '31 days', NOW() - INTERVAL '30 days'),

(1, 4, E'# PetCircle - 以宠物为中心的社交网络\n\n## 产品定位\n不是「养宠物的人的社交」，而是「宠物之间的社交」。\n\n## 差异化策略\n核心差异：**宠物人格化**。每只宠物都有AI生成的性格画像和社交风格。\n\n## MVP 功能清单\n- 宠物个人页（头像、日记、健康档案）\n- 智能配对（根据性格/体型/距离）\n- 宠物朋友圈（纯宠物内容流）\n- 线下聚会组织\n\n## 技术方案\n- Flutter 跨平台客户端\n- Node.js + PostgreSQL 后端\n- 图像识别自动标记宠物品种\n\n## 冷启动策略\n1. 种子用户：本地宠物店合作\n2. 内容激励：每日最佳宠物照片奖励\n3. 线下活动：城市宠物聚会', '[]', 'submitted', 38, NOW() - INTERVAL '29 days', NOW() - INTERVAL '31 days', NOW() - INTERVAL '29 days'),

(1, 7, E'# Pawsitive - 宠物社交解决方案\n\n## 市场分析\n中国宠物市场规模预计 2025 年达 5928 亿元，养宠人群超 7000 万。\n\n## 核心功能矩阵\n### P0（MVP）\n- 宠物主页与动态发布\n- 基于地理位置的宠友发现\n- 即时聊天\n\n### P1\n- 宠物健康记录\n- 服务商入驻\n- 活动组织\n\n## 运营策略\n以城市为单位逐步覆盖，先做一线城市的养猫群体。', '[]', 'submitted', 29, NOW() - INTERVAL '28 days', NOW() - INTERVAL '30 days', NOW() - INTERVAL '28 days'),

(1, 10, E'# 毛球星球 - 产品方案\n\n## 一句话定义\n宠物版的小红书 + 大众点评。\n\n## 功能设计\n1. 内容社区：短视频 + 图文\n2. 服务评价：宠物医院/美容店评分\n3. 交易市场：二手宠物用品\n4. 知识百科：养宠指南\n\n## 商业模式\n- 广告收入\n- 服务商抽佣\n- 电商带货', '[]', 'submitted', 21, NOW() - INTERVAL '27 days', NOW() - INTERVAL '29 days', NOW() - INTERVAL '27 days'),

(1, 14, E'# WoofWoof 产品设计文档\n\n## 用户痛点深挖\n1. 遛狗时想找附近的狗友但没有渠道\n2. 出差时找不到靠谱的宠物寄养\n3. 宠物生病不知道该去哪家医院\n\n## 解决方案\n围绕「宠物生活服务」建立社交关系，不做纯内容社区。\n\n## 技术架构\n- React Native 客户端\n- Go + gRPC 微服务\n- Redis 缓存 + ElasticSearch 搜索\n- 阿里云 OSS 存储', '[]', 'submitted', 17, NOW() - INTERVAL '26 days', NOW() - INTERVAL '28 days', NOW() - INTERVAL '26 days');

-- Idea 2: 独立开发者收入仪表盘 (4 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(2, 1, E'# IndieMetrics - 独立开发者收入聚合\n\n## 问题验证\n调研了 50 位独立开发者，87% 使用 Excel 手动汇总多平台收入。\n\n## 核心功能\n1. **多平台对接**：App Store Connect API、Google Play Console、Stripe、Gumroad、Paddle\n2. **实时仪表盘**：总收入、MRR、ARR、增长趋势\n3. **货币自动换算**：支持 30+ 货币实时汇率\n4. **收入预测**：基于历史数据的趋势预测\n\n## 定价策略\n- Free：3 个数据源，30 天历史\n- Pro $9/月：无限数据源，完整历史，导出功能\n- Team $29/月：多人查看，API 接口', '[]', 'submitted', 55, NOW() - INTERVAL '27 days', NOW() - INTERVAL '28 days', NOW() - INTERVAL '27 days'),

(2, 6, E'# RevenuePulse 方案\n\n## 产品理念\n不做大而全的财务工具，只做「打开就能看到今天赚了多少」这一件事。\n\n## MVP 策略\n第一版只做 Stripe + Gumroad，覆盖 80% 独立开发者的主要收入来源。\n\n## 技术选型\n- Next.js 全栈\n- Vercel 部署\n- OAuth 接入各平台 API\n- Cron job 定时同步数据', '[]', 'submitted', 33, NOW() - INTERVAL '26 days', NOW() - INTERVAL '27 days', NOW() - INTERVAL '26 days'),

(2, 9, E'# MoneyBoard\n\n## 竞品调研\n- Baremetrics：功能强但贵（$129/月起）\n- SimpleAnalytics：只做网站分析\n- 缺口：没有专门给独立开发者的轻量收入看板\n\n## 差异化\n1. 价格亲民（$5/月）\n2. 5 分钟完成所有平台对接\n3. 移动端优先\n\n## 功能清单\n- 收入日历热力图\n- 每日/周/月收入邮件报告\n- 目标追踪（设定月收入目标）\n- 公开收入页面（可选，用于 Build in Public）', '[]', 'submitted', 28, NOW() - INTERVAL '25 days', NOW() - INTERVAL '26 days', NOW() - INTERVAL '25 days'),

(2, 15, E'# 独立开发者收入仪表盘\n\n## 关键洞察\nBuild in Public 趋势下，很多独立开发者愿意公开收入数据。\n\n## 产品 = 工具 + 社区\n- 工具：多平台收入聚合\n- 社区：匿名收入排行榜、里程碑分享\n\n## 增长飞轮\n用户公开收入页 → 其他开发者看到 → 注册使用 → 更多公开页面', '[]', 'submitted', 19, NOW() - INTERVAL '24 days', NOW() - INTERVAL '25 days', NOW() - INTERVAL '24 days');

-- Idea 3: AI 面试模拟器 (5 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(3, 3, E'# MockInterview AI\n\n## 产品愿景\n让每个人都有一个 24/7 随时可用的面试教练。\n\n## 功能设计\n1. **定制面试**：选择公司 + 职位 → 生成针对性问题\n2. **实时模拟**：语音对话，AI 扮演面试官\n3. **即时反馈**：回答评分、改进建议、参考答案\n4. **弱项训练**：根据表现自动推荐练习题\n\n## 技术方案\n- GPT-4 + TTS/STT 实现语音面试\n- 公司题库来自 Glassdoor/面经\n- 评分模型基于 STAR 法则\n\n## 变现\n- 免费：每天 3 题\n- Premium $19/月：无限练习 + 视频回放', '[]', 'submitted', 61, NOW() - INTERVAL '25 days', NOW() - INTERVAL '26 days', NOW() - INTERVAL '25 days'),

(3, 8, E'# InterviewGPT 产品方案\n\n## 市场机会\n全球面试准备市场规模约 $2.4B，在线化趋势明显。\n\n## 核心差异\n不只是问答练习，而是**完整的面试旅程**：简历分析 → 弱点定位 → 定制训练 → 模拟面试 → 面试复盘。\n\n## 分阶段交付\nPhase 1: 文字面试模拟\nPhase 2: 语音面试\nPhase 3: 视频面试 + 肢体语言分析', '[]', 'submitted', 48, NOW() - INTERVAL '24 days', NOW() - INTERVAL '25 days', NOW() - INTERVAL '24 days'),

(3, 11, E'# 面试达人\n\n## 精准定位\n只做技术面试，只做中国市场。\n\n## 题库来源\n- 牛客网面经\n- 力扣讨论区\n- 各大厂面试真题（脱敏处理）\n\n## 功能\n1. 按公司/岗位/轮次筛选题目\n2. AI 模拟不同面试官风格\n3. 代码面试 + 白板协作\n4. 社区互评', '[]', 'submitted', 35, NOW() - INTERVAL '23 days', NOW() - INTERVAL '24 days', NOW() - INTERVAL '23 days'),

(3, 16, E'# PrepAI\n\n## 思路\n把面试准备游戏化，像 Duolingo 学语言一样练面试。\n\n## 核心机制\n- 每日挑战：3 道随机面试题\n- 连续打卡奖励\n- 排行榜（按职位分类）\n- 成就系统\n\n## 技术\n- React + Tailwind 前端\n- Python FastAPI 后端\n- OpenAI API 生成和评估', '[]', 'submitted', 26, NOW() - INTERVAL '22 days', NOW() - INTERVAL '24 days', NOW() - INTERVAL '22 days'),

(3, 20, E'# HireReady\n\n## 独特视角\n不是模拟面试，是**面试策略教练**。帮你分析 JD、研究面试官背景、制定每轮面试策略。\n\n## 功能\n1. JD 深度解析：提取关键技能要求\n2. 面试官 LinkedIn 分析：了解对方技术偏好\n3. 策略建议：每个问题该强调什么\n4. 薪资谈判助手', '[]', 'submitted', 18, NOW() - INTERVAL '21 days', NOW() - INTERVAL '23 days', NOW() - INTERVAL '21 days');

-- Idea 4: 团队知识图谱 (4 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(4, 5, E'# TeamBrain - 团队知识图谱方案\n\n## 问题\n> 新人问："这个服务是谁写的？出了问题找谁？"\n> 答："你去问 Tom，他可能知道。"\n\n## 解决方案\n自动从 Git 历史、Slack 对话、Notion 文档中构建人-知识-项目的关联图谱。\n\n## 核心查询\n- "谁最了解 payment 模块？" → 返回 Git 贡献者排名\n- "关于 Redis 缓存策略的讨论在哪？" → 返回 Slack 对话链接\n- "上次数据库迁移出过什么问题？" → 返回 postmortem 文档', '[]', 'submitted', 44, NOW() - INTERVAL '23 days', NOW() - INTERVAL '24 days', NOW() - INTERVAL '23 days'),

(4, 12, E'# KnowGraph 产品设计\n\n## 切入点\n不做大知识图谱，先做最小有价值的场景：**Who Knows What**。\n\n## MVP\n只对接 GitHub，通过 commit 历史和 code review 记录，自动建立「谁擅长什么」的映射。\n\n## 展示形式\n- Slack Bot：在频道中 @bot "谁了解 GraphQL？"\n- Web 仪表盘：可视化知识分布', '[]', 'submitted', 37, NOW() - INTERVAL '22 days', NOW() - INTERVAL '23 days', NOW() - INTERVAL '22 days'),

(4, 17, E'# 团队大脑\n\n## 技术方案\n- 向量数据库（Qdrant）存储知识嵌入\n- LLM 理解自然语言查询\n- 增量索引，不影响日常工作流\n\n## 隐私设计\n- 所有数据本地处理\n- 支持私有化部署\n- 细粒度权限控制', '[]', 'submitted', 22, NOW() - INTERVAL '21 days', NOW() - INTERVAL '22 days', NOW() - INTERVAL '21 days'),

(4, 21, E'# CollectiveIQ\n\n## 差异化\n不是搜索工具，是**知识发现**工具。主动推送你可能需要但不知道存在的知识。\n\n## 使用场景\n- 写代码时自动推荐相关内部文档\n- 提 PR 时自动找到最合适的 reviewer\n- 新人入职时自动生成个性化学习路径', '[]', 'submitted', 15, NOW() - INTERVAL '20 days', NOW() - INTERVAL '21 days', NOW() - INTERVAL '20 days');

-- Idea 5: 播客笔记助手 (4 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(5, 1, E'# PodNotes 产品方案\n\n## 痛点验证\n"上周某期播客里有人推荐了一本书，但我完全想不起来是哪期了。"\n\n## 核心功能\n1. 自动转录 + 摘要\n2. 关键信息提取（人名、书名、工具、观点）\n3. 时间戳标记\n4. 全文搜索\n\n## 技术\n- Whisper API 转录\n- GPT-4 提取结构化信息\n- 向量搜索支持语义查询', '[]', 'submitted', 39, NOW() - INTERVAL '21 days', NOW() - INTERVAL '22 days', NOW() - INTERVAL '21 days'),

(5, 13, E'# 播客回忆录\n\n## 产品形态\n浏览器插件 + App，在你常用的播客平台上叠加笔记层。\n\n## 功能\n- 收听时一键标记精彩片段\n- AI 自动总结每期核心观点\n- 跨集搜索\n- 知识卡片：自动提取并归类', '[]', 'submitted', 30, NOW() - INTERVAL '20 days', NOW() - INTERVAL '21 days', NOW() - INTERVAL '20 days'),

(5, 18, E'# EarMark\n\n## 独特角度\n不是笔记工具，是**播客内容的 Google**。\n\n## 技术壁垒\n- 自建中英文播客索引（覆盖 10 万+ 节目）\n- 语义搜索：用自然语言搜索播客内容\n- 知识图谱：关联不同播客中提到的相同主题', '[]', 'submitted', 25, NOW() - INTERVAL '19 days', NOW() - INTERVAL '20 days', NOW() - INTERVAL '19 days'),

(5, 22, E'# ListenLearn\n\n## 目标\n把播客从「消遣」变成「学习工具」。\n\n## 功能\n1. AI 生成思维导图\n2. 自动关联推荐相关播客集\n3. 笔记导出到 Notion/Obsidian\n4. 学习进度追踪', '[]', 'submitted', 16, NOW() - INTERVAL '18 days', NOW() - INTERVAL '19 days', NOW() - INTERVAL '18 days');

-- Idea 6: 远程团队异步站会 (3 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(6, 3, E'# AsyncStandup\n\n## 核心理念\n站会的目的是同步信息，不是开会。\n\n## 产品设计\n- 每天下班前花 2 分钟填写：做了什么、计划做什么、是否被阻塞\n- AI 自动生成团队日报\n- 异常检测：连续 3 天同一任务未完成自动提醒\n\n## 与 Slack 深度集成\n- /standup 命令提交更新\n- 自动发送团队摘要到指定频道', '[]', 'submitted', 31, NOW() - INTERVAL '19 days', NOW() - INTERVAL '20 days', NOW() - INTERVAL '19 days'),

(6, 9, E'# DailySync 方案\n\n## 差异化\n不只是文字，支持 60 秒短视频更新。有时候面对面说两句比写一段文字更高效。\n\n## 功能\n1. 文字/视频/语音三种更新方式\n2. AI 自动翻译（解决跨语言团队问题）\n3. 周报自动生成\n4. 与 Jira/Linear 任务联动', '[]', 'submitted', 24, NOW() - INTERVAL '18 days', NOW() - INTERVAL '19 days', NOW() - INTERVAL '18 days'),

(6, 19, E'# TeamPulse\n\n## 核心功能\n- 异步签到\n- 情绪追踪（今天工作状态 1-5 分）\n- 阻塞自动匹配（谁能帮忙解决？）\n- 管理者视角仪表盘', '[]', 'submitted', 14, NOW() - INTERVAL '17 days', NOW() - INTERVAL '18 days', NOW() - INTERVAL '17 days');

-- Idea 7: GitHub Star 追踪器 (4 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(7, 2, E'# StarKeeper\n\n## 问题\nGitHub Star 是开发者的书签，但完全没有管理功能。\n\n## 方案\n1. 自动分类：根据项目 README 和语言自动打标签\n2. 更新追踪：新版本发布、重大变更通知\n3. 死项目检测：超过 6 个月没更新的项目标记\n4. 替代品推荐：某个项目不维护了推荐替代方案\n\n## 变现\n- 免费：100 个 Star 管理\n- Pro $5/月：无限 Star + 高级通知', '[]', 'submitted', 46, NOW() - INTERVAL '17 days', NOW() - INTERVAL '18 days', NOW() - INTERVAL '17 days'),

(7, 8, E'# GitStar Pro\n\n## 浏览器插件方案\n不做独立 App，做 GitHub 增强插件。\n\n## 功能\n- Star 列表增强：搜索、筛选、排序\n- 项目健康度评分\n- Release 通知（浏览器通知）\n- 每周精选邮件', '[]', 'submitted', 34, NOW() - INTERVAL '16 days', NOW() - INTERVAL '17 days', NOW() - INTERVAL '16 days'),

(7, 14, E'# Constellation\n\n## 社交化方向\n看看其他和你 Star 品味相似的开发者都在关注什么。\n\n## 功能\n1. Star 品味匹配\n2. 协作列表（awesome-list 升级版）\n3. 项目讨论社区', '[]', 'submitted', 20, NOW() - INTERVAL '15 days', NOW() - INTERVAL '16 days', NOW() - INTERVAL '15 days'),

(7, 23, E'# StarScope\n\n## 技术方案\n- GitHub API + Webhook\n- 定时任务抓取 Release 信息\n- LLM 总结 changelog\n- 邮件 + Slack 通知', '[]', 'submitted', 12, NOW() - INTERVAL '14 days', NOW() - INTERVAL '15 days', NOW() - INTERVAL '14 days');

-- Idea 8: 个人碳足迹追踪 (3 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(8, 4, E'# CarbonLite 产品方案\n\n## 游戏化设计\n- 每日碳预算（像记账 App）\n- 挑战任务：本周骑车通勤 3 天\n- 排行榜：和朋友比谁更环保\n- 成就徽章\n\n## 数据采集\n- 交通：接入地图 API 自动记录\n- 饮食：拍照识别食物计算碳排\n- 购物：对接电商订单', '[]', 'submitted', 27, NOW() - INTERVAL '15 days', NOW() - INTERVAL '16 days', NOW() - INTERVAL '15 days'),

(8, 11, E'# GreenStep\n\n## 简单至上\n碳计算太复杂，用户不需要精确数据，需要的是方向感。\n\n## 设计\n- 只有三个级别：低碳日/中碳日/高碳日\n- 每天只需回答 3 个问题\n- 月度报告看趋势\n\n## 合作\n- 碳中和品牌合作\n- 种树公益项目', '[]', 'submitted', 19, NOW() - INTERVAL '14 days', NOW() - INTERVAL '15 days', NOW() - INTERVAL '14 days'),

(8, 24, E'# FootprintTracker\n\n## 技术方案\n- 微信小程序（降低使用门槛）\n- 接入高德/百度地图 API\n- 食物碳排数据库\n- 机器学习优化估算模型', '[]', 'submitted', 11, NOW() - INTERVAL '13 days', NOW() - INTERVAL '14 days', NOW() - INTERVAL '13 days');

-- Idea 9-13: 3-4 contributions each (abbreviated for space)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
-- Idea 9: 技术写作 AI 助手 (4 contributions)
(9, 2, E'# WriteDev - 技术写作助手\n\n## 核心功能\n- 代码上下文理解：粘贴代码自动生成文档\n- 多格式输出：博客、API 文档、README\n- 风格迁移：学习你的写作风格\n- SEO 优化：自动添加关键词和结构化数据\n\n## 技术栈\n- VS Code 插件 + Web 编辑器\n- GPT-4 + 代码分析引擎\n- Markdown 实时预览', '[]', 'submitted', 36, NOW() - INTERVAL '13 days', NOW() - INTERVAL '14 days', NOW() - INTERVAL '13 days'),
(9, 7, E'# CodeScribe\n\n## 定位\n不是通用写作助手，是**代码到文档的翻译器**。\n\n## 功能\n1. 一键生成 API 文档\n2. 代码注释自动生成\n3. CHANGELOG 自动更新\n4. 多语言翻译', '[]', 'submitted', 28, NOW() - INTERVAL '12 days', NOW() - INTERVAL '13 days', NOW() - INTERVAL '12 days'),
(9, 15, E'# DocBot\n\n## 差异化\n集成到 CI/CD 流程，每次 PR 合并自动更新文档。\n\n## 使用场景\n- 新模块开发 → 自动生成设计文档骨架\n- API 变更 → 自动更新 API 文档\n- Bug 修复 → 自动生成 postmortem 模板', '[]', 'submitted', 20, NOW() - INTERVAL '11 days', NOW() - INTERVAL '12 days', NOW() - INTERVAL '11 days'),
(9, 25, E'# TechPen\n\n## 产品思路\n让写技术文章像发推一样简单。\n\n## 功能\n- 语音转文章\n- 代码截图美化\n- 一键发布到多平台（掘金/Medium/Dev.to）', '[]', 'submitted', 13, NOW() - INTERVAL '10 days', NOW() - INTERVAL '11 days', NOW() - INTERVAL '10 days'),

-- Idea 10: 开源项目赞助匹配 (3 contributions)
(10, 5, E'# FundOSS\n\n## 双边市场\n- 供给侧：开源维护者注册项目，展示影响力数据\n- 需求侧：企业按技术栈和预算寻找赞助对象\n\n## 匹配算法\n根据企业使用的开源依赖自动推荐赞助对象。\n\n## 变现\n- 企业端收取撮合服务费 10%\n- 提供赞助效果报告', '[]', 'submitted', 30, NOW() - INTERVAL '11 days', NOW() - INTERVAL '12 days', NOW() - INTERVAL '11 days'),
(10, 13, E'# OpenSponsor\n\n## 核心价值\n让开源赞助像买 SaaS 订阅一样简单。\n\n## 功能\n1. 依赖扫描：上传 package.json 看你用了哪些开源项目\n2. 赞助建议：按使用程度推荐赞助金额\n3. 一键赞助：统一支付，自动分配给各项目\n4. 赞助报告：年度开源贡献证书', '[]', 'submitted', 22, NOW() - INTERVAL '10 days', NOW() - INTERVAL '11 days', NOW() - INTERVAL '10 days'),
(10, 26, E'# PatchFund\n\n## 创新点\n不做赞助，做**开源悬赏**。企业发布需要的功能/修复，开源社区认领完成后获得报酬。\n\n## 机制\n- 企业发布 bounty\n- 开发者提交 PR\n- 审核通过后自动打款', '[]', 'submitted', 14, NOW() - INTERVAL '9 days', NOW() - INTERVAL '10 days', NOW() - INTERVAL '9 days'),

-- Idea 11: 城市探索盲盒 (3 contributions)
(11, 6, E'# CityBlind\n\n## 核心体验\n每周收到一个神秘坐标，导航过去发现惊喜。\n\n## 内容来源\n- 本地博主投稿\n- 用户打卡推荐\n- 编辑精选\n\n## 盈利\n- 会员订阅（精选盲盒）\n- 商家合作（探店引流）', '[]', 'submitted', 25, NOW() - INTERVAL '10 days', NOW() - INTERVAL '11 days', NOW() - INTERVAL '10 days'),
(11, 16, E'# Wanderly\n\n## 差异化\n不是推荐热门地点，专门推荐「当地人才知道的地方」。\n\n## 社交元素\n- 探索后发布打卡笔记\n- 区域排行榜：谁探索了最多角落\n- 成就系统：美食猎人/公园达人/胡同专家', '[]', 'submitted', 18, NOW() - INTERVAL '9 days', NOW() - INTERVAL '10 days', NOW() - INTERVAL '9 days'),
(11, 27, E'# UrbanBox\n\n## 技术方案\n- LBS + 兴趣标签匹配\n- UGC + PGC 内容混合\n- 微信小程序 MVP\n- 数据来源：大众点评 API + 自采', '[]', 'submitted', 10, NOW() - INTERVAL '8 days', NOW() - INTERVAL '9 days', NOW() - INTERVAL '8 days'),

-- Idea 12: 代码审查学习平台 (3 contributions)
(12, 3, E'# ReviewSchool\n\n## 产品设计\n- 精选开源项目高质量 PR 案例\n- 按语言/主题/难度分类\n- 用户先自己 review，然后看专家点评\n- 积分和等级系统\n\n## 内容生产\n- 爬取 GitHub 高 Star 项目的 PR\n- AI 筛选有价值的 review 对话\n- 社区贡献 + 专家审核', '[]', 'submitted', 32, NOW() - INTERVAL '9 days', NOW() - INTERVAL '10 days', NOW() - INTERVAL '9 days'),
(12, 10, E'# CodeReviewGym\n\n## 练习模式\n1. 找 Bug 模式：代码中藏了问题，找出来\n2. 改进模式：可以运行但写得不好的代码，提优化建议\n3. 设计模式：评估架构决策\n\n## 分级\nL1 语法问题 → L5 架构级 review', '[]', 'submitted', 23, NOW() - INTERVAL '8 days', NOW() - INTERVAL '9 days', NOW() - INTERVAL '8 days'),
(12, 20, E'# PRMaster\n\n## 思路\n用真实的开源 PR 作为教材，AI 老师逐行讲解为什么这么 review。\n\n## 技术\n- GitHub API 抓取 PR 数据\n- GPT-4 分析 review 逻辑\n- 交互式 diff viewer', '[]', 'submitted', 15, NOW() - INTERVAL '7 days', NOW() - INTERVAL '8 days', NOW() - INTERVAL '7 days'),

-- Idea 13: API 变更通知服务 (3 contributions)
(13, 4, E'# APIWatch\n\n## 监控机制\n1. 定时抓取 API 文档页面 diff\n2. 解析 OpenAPI spec 变更\n3. GitHub Release 监控\n4. 推特/博客关键词监控\n\n## 通知\n- Breaking change → 即时告警\n- Deprecation → 周报\n- 新功能 → 月报\n\n## 迁移辅助\n检测到 breaking change 后，AI 自动生成迁移指南。', '[]', 'submitted', 29, NOW() - INTERVAL '8 days', NOW() - INTERVAL '9 days', NOW() - INTERVAL '8 days'),
(13, 12, E'# ChangeGuard\n\n## CI/CD 集成\n- GitHub Action 每次构建时检查依赖 API 变更\n- 变更评分：影响面（High/Medium/Low）\n- 自动创建 Issue 追踪修复进度\n\n## 覆盖范围\n初期支持 REST API + GraphQL', '[]', 'submitted', 21, NOW() - INTERVAL '7 days', NOW() - INTERVAL '8 days', NOW() - INTERVAL '7 days'),
(13, 28, E'# APIDrift\n\n## 独特方案\n不监控文档，监控**实际 API 行为**。定期发送测试请求，对比响应结构变化。\n\n## 优势\n- 发现文档没写但实际已变更的改动\n- 比爬文档更可靠', '[]', 'submitted', 11, NOW() - INTERVAL '6 days', NOW() - INTERVAL '7 days', NOW() - INTERVAL '6 days');

-- Open ideas contributions (mix of draft and submitted)
-- Idea 14: AI Commit Message (4 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(14, 1, E'# CommitCraft\n\n## 方案概要\n- 分析 git diff 语义\n- 遵循 Conventional Commits\n- 支持中英文\n- VS Code / CLI 双入口', '[]', 'submitted', 8, NOW() - INTERVAL '1 day', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day'),
(14, 5, E'# GitMsg AI\n\n## 差异化\n学习项目历史 commit 风格，生成风格一致的 message。', '[]', 'submitted', 5, NOW() - INTERVAL '12 hours', NOW() - INTERVAL '1 day', NOW() - INTERVAL '12 hours'),
(14, 18, E'# 正在完善中...', '[]', 'draft', 0, NULL, NOW() - INTERVAL '1 day', NOW() - INTERVAL '6 hours'),
(14, 22, E'# AutoCommit\n\n## 功能\n- git hook 自动触发\n- 多种 commit 风格模板\n- 团队统一规范配置', '[]', 'submitted', 3, NOW() - INTERVAL '6 hours', NOW() - INTERVAL '1 day', NOW() - INTERVAL '6 hours');

-- Idea 15: 开发者人体工学 (3 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(15, 7, E'# DevHealth\n\n## 方案\n- 番茄钟 + 姿态提醒结合\n- 每 45 分钟提醒做拉伸\n- 每 2 小时建议眼保健操\n- macOS 菜单栏常驻', '[]', 'submitted', 6, NOW() - INTERVAL '2 days', NOW() - INTERVAL '3 days', NOW() - INTERVAL '2 days'),
(15, 11, E'# PostureGuard\n\n## 创新点\n不用摄像头，通过键盘打字节奏和鼠标移动模式推断疲劳程度。', '[]', 'submitted', 4, NOW() - INTERVAL '1 day', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day'),
(15, 29, E'# 草稿中\n\n准备从 Apple Watch 健康数据切入...', '[]', 'draft', 0, NULL, NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day');

-- Idea 16: 技术播客推荐 (2 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(16, 9, E'# PodMatch\n\n## 推荐算法\n- 基于技术栈标签匹配\n- 协同过滤：相似开发者听什么\n- 内容分析：NLP 提取单集主题', '[]', 'submitted', 4, NOW() - INTERVAL '12 hours', NOW() - INTERVAL '1 day', NOW() - INTERVAL '12 hours'),
(16, 24, E'# TechCast\n\n## MVP\n先做一个精选列表，人工运营 + AI 推荐结合。', '[]', 'draft', 0, NULL, NOW() - INTERVAL '18 hours', NOW() - INTERVAL '6 hours');

-- Idea 17: 开源许可证检查器 (3 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(17, 3, E'# LicenseGuard\n\n## CI 集成方案\n- GitHub Action 一键接入\n- 支持 npm/pip/go mod/cargo\n- 生成合规报告\n- 许可证冲突自动检测', '[]', 'submitted', 7, NOW() - INTERVAL '3 days', NOW() - INTERVAL '4 days', NOW() - INTERVAL '3 days'),
(17, 10, E'# CompliBot\n\n## 差异化\n不只检测，还提供修复建议和替代依赖推荐。', '[]', 'submitted', 5, NOW() - INTERVAL '2 days', NOW() - INTERVAL '3 days', NOW() - INTERVAL '2 days'),
(17, 30, E'# 许可证雷达\n\n## 功能\n- 依赖树可视化\n- 许可证传染性分析\n- SBOM 生成', '[]', 'submitted', 3, NOW() - INTERVAL '1 day', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day');

-- Idea 18: 代码片段搜索 (2 contributions)
INSERT INTO contributions (idea_id, author_id, content, decision_log, status, view_count, submitted_at, created_at, updated_at) VALUES
(18, 6, E'# SnippetSearch\n\n## 方案\n- 本地索引所有项目\n- 向量化代码嵌入\n- 自然语言搜索', '[]', 'submitted', 5, NOW() - INTERVAL '1 day', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day'),
(18, 15, E'# CodeVault 草稿', '[]', 'draft', 0, NULL, NOW() - INTERVAL '1 day', NOW() - INTERVAL '6 hours');


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
