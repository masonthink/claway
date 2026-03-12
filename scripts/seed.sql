-- Seed data for Claway demo
BEGIN;

-- ============================================================
-- Users (5 users with different roles)
-- ============================================================
INSERT INTO users (id, openclaw_id, username, display_name, avatar_url, credits_balance, created_at) VALUES
(1, '', 'alice',  'Alice Chen',   '', 1520.0000, '2026-02-10 08:00:00+08'),
(2, '', 'bob',    'Bob Zhang',    '', 3200.5000, '2026-02-12 10:30:00+08'),
(3, '', 'carol',  'Carol Wang',   '',  860.0000, '2026-02-15 14:00:00+08'),
(4, '', 'david',  'David Li',     '',  450.0000, '2026-02-20 09:00:00+08'),
(5, '', 'eve',    'Eve Liu',      '',    0.0000, '2026-03-01 16:00:00+08');

SELECT setval('users_id_seq', 5);

-- ============================================================
-- Idea 1: AI 笔记助手 (standard, active, 3/9 tasks done)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(1,
 'AI 笔记助手 — NoteFlow',
 '为独立开发者和学生设计的智能笔记工具。支持语音转文字、自动分类、跨笔记关联搜索和 AI 摘要生成。目标是让碎片化知识变成可检索、可复用的知识库。',
 '独立开发者、大学生、知识工作者',
 '现有笔记工具（Notion、Obsidian）学习成本高，缺乏智能整理能力。用户记了大量笔记但很少回顾，知识利用率低。',
 1, 20.00, 'standard', 'active', '2026-02-18 10:00:00+08');

-- Tasks for Idea 1
INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, claimed_by, claimed_at, submitted_at, approved_at, output_content, quality_score, cost_usd_accumulated) VALUES
(1, 1, 'D1', '竞品分析报告',
 'Research and analyze the competitive landscape for AI-powered note-taking tools.',
 '>=3 direct competitors + >=2 indirect competitors analyzed',
 '', 80000, 'approved', 2, '2026-02-19 09:00:00+08', '2026-02-20 18:00:00+08', '2026-02-21 10:00:00+08',
 'Completed competitive analysis covering Notion AI, Mem, Reflect, Obsidian, and Apple Notes.',
 1.20, 0.850000),

(2, 1, 'D2', '目标用户画像',
 'Define target user personas with core pain points and usage scenarios.',
 '2-3 user personas with narrative scenarios',
 '', 60000, 'approved', 3, '2026-02-20 14:00:00+08', '2026-02-22 11:00:00+08', '2026-02-23 09:00:00+08',
 'Created 3 detailed personas: indie developer, graduate student, product manager.',
 1.50, 0.620000),

(3, 1, 'D3', '用户旅程地图',
 'Create user journey maps for core usage scenarios.',
 'User story format with acceptance criteria, P0 features <=10',
 'D1,D2', 100000, 'approved', 2, '2026-02-24 10:00:00+08', '2026-02-27 16:00:00+08', '2026-02-28 11:00:00+08',
 'Mapped 4 core journeys: capture, organize, retrieve, share.',
 1.20, 1.230000),

(4, 1, 'D4', '功能需求文档',
 'Define functional requirements with P0/P1 priority.',
 'Revenue model + pricing rationale + cold-start path',
 'D1,D2', 60000, 'submitted', 4, '2026-03-01 09:00:00+08', '2026-03-05 17:00:00+08', NULL,
 'Defined freemium model with 3 tiers.', NULL, 0.480000),

(5, 1, 'D5', '信息架构图',
 'Define success metrics and information architecture.',
 '1 north star metric + 3-5 process metrics',
 'D1,D2', 40000, 'claimed', 3, '2026-03-02 10:00:00+08', NULL, NULL,
 NULL, NULL, 0.150000),

(6, 1, 'D6', '页面流程图',
 'Design page structure and navigation system.',
 'Complete page list + navigation structure + permission matrix',
 'D3', 60000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0),

(7, 1, 'D7', '交互设计规范',
 'Map core user flows for all P0 features.',
 'Normal path + >=2 exception paths per P0 feature',
 'D3', 80000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0),

(8, 1, 'D8', '视觉设计规范',
 'Create design specifications with component library.',
 'Component library + color palette + typography + spacing',
 'D3', 60000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0),

(9, 1, 'D9', '技术可行性评估',
 'Evaluate tech stack and key risk points.',
 'Tech stack recommendations + risk points + feasibility conclusion',
 'D3', 60000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0);

-- ============================================================
-- Idea 2: 社区团购小程序 (light, active, 1/5 tasks done)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(2,
 '社区团购小程序 — 邻里拼',
 '面向小区居民的社区团购平台。团长发起拼团，邻居一键参团，支持自提点取货。解决生鲜水果最后一公里配送成本高的问题。',
 '社区团长、小区居民、生鲜供应商',
 '现有社区团购平台（美团优选、多多买菜）佣金高、品类受限。小区自发团购用微信群管理混乱，对账困难。',
 2, 15.00, 'light', 'active', '2026-02-25 14:00:00+08');

INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, claimed_by, claimed_at, submitted_at, approved_at, output_content, quality_score, cost_usd_accumulated) VALUES
(10, 2, 'D1', '竞品分析报告',
 'Analyze community group-buying platforms.',
 '>=3 direct competitors analyzed',
 '', 80000, 'approved', 4, '2026-02-26 10:00:00+08', '2026-02-28 15:00:00+08', '2026-03-01 09:00:00+08',
 'Analyzed 美团优选, 多多买菜, 淘菜菜, 兴盛优选, and 3 indie WeChat mini-programs.',
 1.00, 0.720000),

(11, 2, 'D2', '目标用户画像',
 'Define personas for group-buying participants.',
 '2-3 user personas with scenarios',
 '', 60000, 'claimed', 5, '2026-03-02 11:00:00+08', NULL, NULL,
 NULL, NULL, 0.080000),

(12, 2, 'D3', '用户旅程地图',
 'Map user journeys for group leaders and participants.',
 'User stories with acceptance criteria',
 'D1,D2', 100000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0),

(13, 2, 'D4', '功能需求文档',
 'Define functional requirements for the mini-program.',
 'Revenue model + pricing + cold-start path',
 'D1,D2', 60000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0),

(14, 2, 'D5', '信息架构图',
 'Define metrics and information architecture.',
 'North star metric + process metrics',
 'D1,D2', 40000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0);

-- ============================================================
-- Idea 3: 宠物健康管理 App (standard, completed)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(3,
 '宠物健康管理 App — PetCare+',
 '帮助宠物主人记录宠物的饮食、运动、疫苗、体检等健康数据，AI 分析异常指标并提醒就医。支持多宠物管理和兽医在线问诊。',
 '养猫养狗的年轻人、宠物店主、兽医',
 '宠物就医成本高且信息不对称。主人难以判断宠物是否生病，错过最佳治疗时机。现有宠物 App 功能碎片化，缺乏健康数据追踪。',
 3, 25.00, 'standard', 'completed', '2026-01-15 09:00:00+08');

INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, claimed_by, claimed_at, submitted_at, approved_at, quality_score, cost_usd_accumulated) VALUES
(15, 3, 'D1', '竞品分析报告', 'Analyze pet health management apps.', '>=3 competitors', '', 80000, 'approved', 2, '2026-01-16 10:00:00+08', '2026-01-18 16:00:00+08', '2026-01-19 09:00:00+08', 1.20, 0.900000),
(16, 3, 'D2', '目标用户画像', 'Define pet owner personas.', '2-3 personas', '', 60000, 'approved', 4, '2026-01-17 14:00:00+08', '2026-01-20 11:00:00+08', '2026-01-21 10:00:00+08', 1.50, 0.550000),
(17, 3, 'D3', '用户旅程地图', 'Map pet care user journeys.', 'P0 features <=10', 'D1,D2', 100000, 'approved', 2, '2026-01-22 09:00:00+08', '2026-01-26 17:00:00+08', '2026-01-27 11:00:00+08', 1.20, 1.100000),
(18, 3, 'D4', '功能需求文档', 'Define functional requirements.', 'Revenue model', 'D1,D2', 60000, 'approved', 3, '2026-01-23 10:00:00+08', '2026-01-27 15:00:00+08', '2026-01-28 09:00:00+08', 1.00, 0.650000),
(19, 3, 'D5', '信息架构图', 'Define metrics and IA.', 'North star + process metrics', 'D1,D2', 40000, 'approved', 4, '2026-01-24 11:00:00+08', '2026-01-27 18:00:00+08', '2026-01-28 10:00:00+08', 1.20, 0.420000),
(20, 3, 'D6', '页面流程图', 'Design page flows.', 'Complete page list', 'D3', 60000, 'approved', 2, '2026-01-29 09:00:00+08', '2026-02-01 16:00:00+08', '2026-02-02 09:00:00+08', 1.00, 0.580000),
(21, 3, 'D7', '交互设计规范', 'Map interaction flows.', 'P0 flows covered', 'D3', 80000, 'approved', 3, '2026-01-30 10:00:00+08', '2026-02-03 17:00:00+08', '2026-02-04 10:00:00+08', 1.50, 0.870000),
(22, 3, 'D8', '视觉设计规范', 'Create design specs.', 'Component library + colors', 'D3', 60000, 'approved', 5, '2026-02-01 14:00:00+08', '2026-02-05 15:00:00+08', '2026-02-06 09:00:00+08', 1.20, 0.620000),
(23, 3, 'D9', '技术可行性评估', 'Evaluate tech stack.', 'Tech recommendations + risks', 'D3', 60000, 'approved', 4, '2026-02-02 10:00:00+08', '2026-02-06 18:00:00+08', '2026-02-07 11:00:00+08', 1.00, 0.510000);

-- ============================================================
-- Idea 4: 自由职业者接单平台 (standard, active, fresh)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(4,
 '自由职业者接单平台 — FreelanceHub',
 '连接自由职业者和企业的智能匹配平台。基于技能图谱和项目需求自动推荐合适的自由职业者，支持里程碑付款和在线协作。',
 '设计师、开发者、文案、翻译等自由职业者；中小企业、创业团队',
 '现有平台（猪八戒、Fiverr）匹配效率低、抽佣高。自由职业者获客难，企业找到合适人才耗时长。',
 4, 18.00, 'standard', 'active', '2026-03-08 11:00:00+08');

INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, cost_usd_accumulated) VALUES
(24, 4, 'D1', '竞品分析报告', 'Analyze freelance platforms.', '>=3 competitors', '', 80000, 'open', 0),
(25, 4, 'D2', '目标用户画像', 'Define freelancer and client personas.', '2-3 personas', '', 60000, 'open', 0),
(26, 4, 'D3', '用户旅程地图', 'Map matching and hiring journeys.', 'P0 features <=10', 'D1,D2', 100000, 'open', 0),
(27, 4, 'D4', '功能需求文档', 'Define platform requirements.', 'Revenue model', 'D1,D2', 60000, 'open', 0),
(28, 4, 'D5', '信息架构图', 'Define metrics and IA.', 'North star + process metrics', 'D1,D2', 40000, 'open', 0),
(29, 4, 'D6', '页面流程图', 'Design page flows.', 'Complete page list', 'D3', 60000, 'open', 0),
(30, 4, 'D7', '交互设计规范', 'Map interaction flows.', 'P0 flows covered', 'D3', 80000, 'open', 0),
(31, 4, 'D8', '视觉设计规范', 'Create design specs.', 'Component library', 'D3', 60000, 'open', 0),
(32, 4, 'D9', '技术可行性评估', 'Evaluate tech feasibility.', 'Tech recommendations', 'D3', 60000, 'open', 0);

-- ============================================================
-- Idea 5: 商标交易平台 (standard, active, 2/9 tasks done)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(5,
 '商标交易平台 — 标易通',
 '面向中小企业和个体商户的商标交易撮合平台。整合闲置商标资源，通过 AI 估价、智能匹配和在线过户一站式服务，降低商标获取成本和时间。支持商标买卖、授权许可和质押融资。',
 '创业公司创始人、中小企业主、知识产权代理机构、个体工商户',
 '商标注册周期长（6-12个月）且通过率低（约40%）。现有商标交易平台（如中华商标超市、八戒知产）信息不透明、估价混乱、过户流程复杂。中小企业急需商标时缺乏高效获取渠道。',
 1, 22.00, 'standard', 'active', '2026-03-05 09:00:00+08');

-- Tasks for Idea 5
INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, claimed_by, claimed_at, submitted_at, approved_at, output_content, quality_score, cost_usd_accumulated) VALUES
(33, 5, 'D1', '竞品分析报告',
 'Research and analyze existing trademark trading platforms, IP marketplaces, and related services.',
 '>=3 direct competitors + >=2 indirect competitors analyzed',
 '', 80000, 'approved', 3, '2026-03-06 10:00:00+08', '2026-03-08 16:00:00+08', '2026-03-09 09:00:00+08',
 'Analyzed 中华商标超市, 八戒知产, 权大师, 标天下, and indirect competitors like 企查查/天眼查 trademark modules.',
 1.20, 0.880000),

(34, 5, 'D2', '目标用户画像',
 'Define target user personas for trademark buyers, sellers, and IP agencies.',
 '2-3 user personas with narrative scenarios',
 '', 60000, 'approved', 4, '2026-03-07 14:00:00+08', '2026-03-09 11:00:00+08', '2026-03-10 10:00:00+08',
 'Created 3 personas: startup founder needing brand protection, retiring merchant selling unused marks, IP agency managing portfolio.',
 1.50, 0.590000),

(35, 5, 'D3', '用户旅程地图',
 'Map user journeys for trademark buying, selling, and licensing scenarios.',
 'User story format with acceptance criteria, P0 features <=10',
 'D1,D2', 100000, 'claimed', 2, '2026-03-10 09:00:00+08', NULL, NULL,
 NULL, NULL, 0.320000),

(36, 5, 'D4', '功能需求文档',
 'Define functional requirements including AI valuation, smart matching, and online transfer.',
 'Revenue model + pricing rationale + cold-start path',
 'D1,D2', 60000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0),

(37, 5, 'D5', '信息架构图',
 'Define information architecture for trademark catalog, search, and transaction flows.',
 '1 north star metric + 3-5 process metrics',
 'D1,D2', 40000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0),

(38, 5, 'D6', '页面流程图',
 'Design page structure for marketplace, listing, and transfer management.',
 'Complete page list + navigation structure + permission matrix',
 'D3', 60000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0),

(39, 5, 'D7', '交互设计规范',
 'Map interaction flows for search, purchase, and transfer processes.',
 'Normal path + >=2 exception paths per P0 feature',
 'D3', 80000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0),

(40, 5, 'D8', '视觉设计规范',
 'Create design specs for a professional, trust-oriented trading platform.',
 'Component library + color palette + typography + spacing',
 'D3', 60000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0),

(41, 5, 'D9', '技术可行性评估',
 'Evaluate tech stack for trademark data integration, AI valuation model, and secure transactions.',
 'Tech stack recommendations + risk points + feasibility conclusion',
 'D3', 60000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0);

SELECT setval('ideas_id_seq', 5);
SELECT setval('tasks_id_seq', 41);

-- ============================================================
-- Documents for all tasks
-- ============================================================

-- Idea 1 documents (approved tasks have content)
INSERT INTO documents (id, task_id, content, current_version) VALUES
(1,  1,  E'# 竞品分析报告 — AI 笔记助手\n\n## 1. 直接竞品\n\n### 1.1 Notion AI\n- **定位**: All-in-one workspace with AI assistant\n- **优势**: 生态完善，模板丰富，协作能力强\n- **劣势**: 学习曲线陡，AI 功能偏通用，不专注笔记场景\n- **定价**: $10/月（Plus），AI 附加 $8/月\n\n### 1.2 Mem\n- **定位**: AI-first note-taking for professionals\n- **优势**: 自动关联、智能搜索、零组织成本\n- **劣势**: 生态封闭，不支持本地存储，团队功能弱\n- **定价**: $15/月\n\n### 1.3 Reflect\n- **定位**: 端到端加密的 AI 笔记\n- **优势**: 隐私保护强，AI 摘要质量高\n- **劣势**: 功能简单，无协作，市场份额小\n\n## 2. 间接竞品\n\n### 2.1 Obsidian\n- 本地优先、插件生态强，但 AI 能力需第三方插件\n\n### 2.2 Apple Notes\n- 零学习成本，但智能化程度低\n\n## 3. 差异化空间\n- **语音转笔记**: 会议/灵感场景，竞品普遍弱\n- **知识图谱可视化**: 超越简单双链\n- **中文语境优化**: 国际竞品中文支持差', 2),

(2,  2,  E'# 目标用户画像\n\n## Persona 1: 独立开发者 Alex\n- **年龄**: 28岁\n- **痛点**: 学习新技术时笔记分散在多个平台，无法快速找到之前记录的解决方案\n- **场景**: 边写代码边记 debug 日志，需要快速语音输入\n- **预算**: 愿付 ¥30-50/月\n\n## Persona 2: 研究生 小林\n- **年龄**: 24岁\n- **痛点**: 论文阅读笔记多但缺乏系统整理，写综述时找不到引用\n- **场景**: 课堂录音转文字，文献笔记自动关联\n- **预算**: 愿付 ¥15-25/月\n\n## Persona 3: 产品经理 Jenny\n- **年龄**: 32岁\n- **痛点**: 需求调研、用户访谈、竞品分析散落各处\n- **场景**: 会议纪要自动生成 action items，跨项目知识复用\n- **预算**: 公司报销，¥50-100/月', 2),

(3,  3,  E'# 用户旅程地图\n\n## Journey 1: 知识捕获\n**触发** → 灵感闪现/会议开始 → 打开 App → 语音/文字输入 → AI 自动分类 → 确认保存\n\n## Journey 2: 知识组织\n**触发** → 周末整理 → 查看本周笔记 → AI 建议关联 → 确认/调整标签 → 生成知识图谱\n\n## Journey 3: 知识检索\n**触发** → 遇到问题 → 自然语言搜索 → AI 返回相关笔记 + 摘要 → 直接使用\n\n## Journey 4: 知识分享\n**触发** → 写文章/汇报 → 选择笔记 → AI 生成大纲 → 导出分享\n\n## P0 Features (8)\n1. 文字笔记创建与编辑\n2. 语音转文字\n3. AI 自动分类\n4. 自然语言搜索\n5. AI 摘要生成\n6. 笔记关联推荐\n7. 标签管理\n8. 导出（Markdown/PDF）', 2),

(4,  4,  '', 1),
(5,  5,  '', 1),
(6,  6,  '', 1),
(7,  7,  '', 1),
(8,  8,  '', 1),
(9,  9,  '', 1);

-- Idea 2 documents
INSERT INTO documents (id, task_id, content, current_version) VALUES
(10, 10, E'# 竞品分析 — 社区团购\n\n## 直接竞品\n1. **美团优选** — 流量大但佣金 15-20%，品类标准化\n2. **多多买菜** — 低价策略，依赖拼多多流量\n3. **淘菜菜** — 阿里系供应链，覆盖有限\n4. **兴盛优选** — 下沉市场强，湖南起家\n\n## 间接竞品\n- 微信群手动团购\n- 小区物业自营团购\n\n## 差异化\n- 去中心化：团长自主定价\n- 本地化供应链对接\n- 低佣金（5%）吸引团长', 1),
(11, 11, '', 1),
(12, 12, '', 1),
(13, 13, '', 1),
(14, 14, '', 1);

-- Idea 3 documents (all approved, brief content)
INSERT INTO documents (id, task_id, content, current_version) VALUES
(15, 15, '# 竞品分析 — 宠物健康管理', 1),
(16, 16, '# 目标用户画像', 1),
(17, 17, '# 用户旅程地图', 1),
(18, 18, '# 功能需求文档', 1),
(19, 19, '# 信息架构图', 1),
(20, 20, '# 页面流程图', 1),
(21, 21, '# 交互设计规范', 1),
(22, 22, '# 视觉设计规范', 1),
(23, 23, '# 技术可行性评估', 1);

-- Idea 4 documents (all empty)
INSERT INTO documents (id, task_id, content, current_version) VALUES
(24, 24, '', 1),
(25, 25, '', 1),
(26, 26, '', 1),
(27, 27, '', 1),
(28, 28, '', 1),
(29, 29, '', 1),
(30, 30, '', 1),
(31, 31, '', 1),
(32, 32, '', 1);

-- Idea 5 documents
INSERT INTO documents (id, task_id, content, current_version) VALUES
(33, 33, E'# 竞品分析报告 — 商标交易平台\n\n## 1. 直接竞品\n\n### 1.1 中华商标超市网\n- **定位**: 国内最大商标转让平台\n- **优势**: 商标库量大（200万+），品牌知名度高\n- **劣势**: 界面老旧，搜索体验差，估价不透明\n- **佣金**: 交易额 5-10%\n\n### 1.2 八戒知产（猪八戒旗下）\n- **定位**: 一站式知识产权服务\n- **优势**: 流量大，服务链完整（注册+转让+维权）\n- **劣势**: 商标交易非核心业务，匹配效率低\n\n### 1.3 权大师\n- **定位**: 智能知识产权服务平台\n- **优势**: AI 辅助检索，数据可视化好\n- **劣势**: 侧重注册，交易功能弱\n\n### 1.4 标天下\n- **定位**: 商标交易垂直平台\n- **优势**: 专注交易，流程清晰\n- **劣势**: 规模小，商标库有限\n\n## 2. 间接竞品\n\n### 2.1 企查查/天眼查 商标模块\n- 查询功能强但无交易能力\n\n### 2.2 知识产权代理机构（线下）\n- 信任度高但效率低、费用高\n\n## 3. 差异化空间\n- **AI 估价模型**: 基于商标类别、注册年限、行业热度、近似商标成交价智能定价\n- **智能匹配**: 买家需求画像 × 商标特征向量，主动推荐\n- **过户全托管**: 在线签约 + 材料代办 + 进度追踪，交易周期压缩至 2-3 个月', 2),

(34, 34, E'# 目标用户画像 — 商标交易\n\n## Persona 1: 创业公司创始人 李明\n- **年龄**: 32岁\n- **背景**: 互联网创业者，正在筹备新消费品牌\n- **痛点**: 注册商标周期太长（8个月+），品牌上线时间紧迫；不确定想要的名字是否能注册成功\n- **场景**: 直接购买已注册商标，缩短品牌上线时间\n- **预算**: 5000-30000元\n\n## Persona 2: 个体商户 王姐\n- **年龄**: 45岁\n- **背景**: 经营了 10 年的服装店，注册了 3 个商标但只在用 1 个\n- **痛点**: 闲置商标每年要交维护费，但不知道值多少钱、怎么卖\n- **场景**: 挂牌出售闲置商标，变现知识产权\n- **预算**: 希望零成本挂牌，成交后抽佣\n\n## Persona 3: 知识产权代理人 张律师\n- **年龄**: 38岁\n- **背景**: 运营一家小型知识产权代理所，管理 500+ 客户商标\n- **痛点**: 客户经常问"这个商标值多少钱""帮我找个XX类的商标"，手动匹配效率低\n- **场景**: 批量管理客户商标资产，用平台工具辅助估价和匹配\n- **预算**: 愿付 ¥200-500/月 SaaS 费用', 2),

(35, 35, '', 1),
(36, 36, '', 1),
(37, 37, '', 1),
(38, 38, '', 1),
(39, 39, '', 1),
(40, 40, '', 1),
(41, 41, '', 1);

SELECT setval('documents_id_seq', 41);

-- ============================================================
-- Token usage logs (for tasks with compute cost)
-- ============================================================
INSERT INTO token_usage_logs (user_id, task_id, model, tokens_in, tokens_out, cost_usd, timestamp) VALUES
-- Idea 1
(2, 1, 'claude-3.5-sonnet', 45000, 12000, 0.420000, '2026-02-20 14:00:00+08'),
(2, 1, 'claude-3.5-sonnet', 38000, 15000, 0.430000, '2026-02-20 17:00:00+08'),
(3, 2, 'claude-3.5-sonnet', 32000, 18000, 0.620000, '2026-02-22 10:00:00+08'),
(2, 3, 'claude-3.5-sonnet', 55000, 20000, 0.650000, '2026-02-26 11:00:00+08'),
(2, 3, 'claude-3.5-sonnet', 48000, 18000, 0.580000, '2026-02-27 15:00:00+08'),
(4, 4, 'claude-3.5-sonnet', 28000, 14000, 0.480000, '2026-03-05 15:00:00+08'),
(3, 5, 'claude-3.5-sonnet', 12000, 5000,  0.150000, '2026-03-03 14:00:00+08'),
-- Idea 2
(4, 10, 'claude-3.5-sonnet', 40000, 16000, 0.720000, '2026-02-28 14:00:00+08'),
(5, 11, 'claude-3.5-sonnet', 8000,  3000,  0.080000, '2026-03-03 11:00:00+08'),
-- Idea 3
(2, 15, 'claude-3.5-sonnet', 50000, 14000, 0.900000, '2026-01-18 15:00:00+08'),
(4, 16, 'claude-3.5-sonnet', 30000, 16000, 0.550000, '2026-01-20 10:00:00+08'),
(2, 17, 'claude-3.5-sonnet', 60000, 22000, 1.100000, '2026-01-26 16:00:00+08'),
(3, 18, 'claude-3.5-sonnet', 35000, 15000, 0.650000, '2026-01-27 14:00:00+08'),
(4, 19, 'claude-3.5-sonnet', 24000, 10000, 0.420000, '2026-01-27 17:00:00+08'),
(2, 20, 'claude-3.5-sonnet', 32000, 14000, 0.580000, '2026-02-01 15:00:00+08'),
(3, 21, 'claude-3.5-sonnet', 48000, 20000, 0.870000, '2026-02-03 16:00:00+08'),
(5, 22, 'claude-3.5-sonnet', 34000, 15000, 0.620000, '2026-02-05 14:00:00+08'),
(4, 23, 'claude-3.5-sonnet', 28000, 12000, 0.510000, '2026-02-06 17:00:00+08'),
-- Idea 5
(3, 33, 'claude-3.5-sonnet', 48000, 15000, 0.880000, '2026-03-08 15:00:00+08'),
(4, 34, 'claude-3.5-sonnet', 30000, 16000, 0.590000, '2026-03-09 10:00:00+08'),
(2, 35, 'claude-3.5-sonnet', 18000, 8000,  0.320000, '2026-03-11 14:00:00+08');

-- ============================================================
-- Contributions (for approved tasks)
-- ============================================================
-- Idea 1 approved tasks
INSERT INTO contributions (idea_id, task_id, user_id, cost_usd, quality_score, weighted_score, weight_percent) VALUES
(1, 1, 2, 0.850000, 1.20, 1.020000, 0),
(1, 2, 3, 0.620000, 1.50, 0.930000, 0),
(1, 3, 2, 1.230000, 1.20, 1.476000, 0);

-- Idea 2 approved tasks
INSERT INTO contributions (idea_id, task_id, user_id, cost_usd, quality_score, weighted_score, weight_percent) VALUES
(2, 10, 4, 0.720000, 1.00, 0.720000, 100.0000);

-- Idea 5 approved tasks
INSERT INTO contributions (idea_id, task_id, user_id, cost_usd, quality_score, weighted_score, weight_percent) VALUES
(5, 33, 3, 0.880000, 1.20, 1.056000, 0),
(5, 34, 4, 0.590000, 1.50, 0.885000, 0);

-- Idea 3 all approved
INSERT INTO contributions (idea_id, task_id, user_id, cost_usd, quality_score, weighted_score, weight_percent) VALUES
(3, 15, 2, 0.900000, 1.20, 1.080000, 15.2000),
(3, 16, 4, 0.550000, 1.50, 0.825000, 11.6100),
(3, 17, 2, 1.100000, 1.20, 1.320000, 18.5900),
(3, 18, 3, 0.650000, 1.00, 0.650000, 9.1500),
(3, 19, 4, 0.420000, 1.20, 0.504000, 7.0900),
(3, 20, 2, 0.580000, 1.00, 0.580000, 8.1600),
(3, 21, 3, 0.870000, 1.50, 1.305000, 18.3700),
(3, 22, 5, 0.620000, 1.20, 0.744000, 10.4700),
(3, 23, 4, 0.510000, 1.00, 0.510000, 7.1800);

-- ============================================================
-- PRD for completed idea 3
-- ============================================================
INSERT INTO prds (id, idea_id, content, published_at, price_credits, read_count) VALUES
(1, 3,
 E'# PetCare+ — 宠物健康管理 App PRD\n\n## 产品概述\n面向年轻宠物主人的一站式健康管理工具...\n\n## 核心功能\n1. 健康数据记录（体重、饮食、运动）\n2. 疫苗与体检提醒\n3. AI 异常指标预警\n4. 兽医在线问诊\n5. 多宠物管理\n\n## 商业模式\n- 基础版免费\n- Pro 版 ¥25/月（AI 分析 + 兽医问诊）\n- 宠物医院 B 端合作',
 '2026-02-10 09:00:00+08', 500.0000, 3);

SELECT setval('prds_id_seq', 1);

-- ============================================================
-- Credit transactions
-- ============================================================
INSERT INTO credit_transactions (user_id, type, amount, reference_type, reference_id, description, created_at) VALUES
-- Bob (user 2) earned from contributions
(2, 'earn_contribute', 1020.0000, 'contribution', 1, 'D1 竞品分析 — AI 笔记助手', '2026-02-21 10:00:00+08'),
(2, 'earn_contribute', 1476.0000, 'contribution', 3, 'D3 用户旅程 — AI 笔记助手', '2026-02-28 11:00:00+08'),
-- Carol (user 3) earned
(3, 'earn_contribute', 930.0000, 'contribution', 2, 'D2 用户画像 — AI 笔记助手', '2026-02-23 09:00:00+08'),
-- David (user 4) earned
(4, 'earn_contribute', 720.0000, 'contribution', 4, 'D1 竞品分析 — 社区团购', '2026-03-01 09:00:00+08'),
-- PRD purchase
(5, 'spend_read', -500.0000, 'prd', 1, '购买 PetCare+ PRD', '2026-03-05 14:00:00+08');

-- ============================================================
-- OAuth accounts (X login)
-- ============================================================
INSERT INTO user_oauth_accounts (user_id, provider, provider_user_id, provider_username) VALUES
(1, 'x', '1001', 'alice_chen'),
(2, 'x', '1002', 'bob_builds'),
(3, 'x', '1003', 'carol_designs'),
(4, 'x', '1004', 'david_codes'),
(5, 'x', '1005', 'eve_creates');

COMMIT;
