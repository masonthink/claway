-- Seed data for Claway demo (4 tasks per idea: D1-D4)
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
-- Idea 1: AI 笔记助手 (active, D1+D2 approved, D3 submitted, D4 open)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(1,
 'AI 笔记助手 — NoteFlow',
 '为独立开发者和学生设计的智能笔记工具。支持语音转文字、自动分类、跨笔记关联搜索和 AI 摘要生成。目标是让碎片化知识变成可检索、可复用的知识库。',
 '独立开发者、大学生、知识工作者',
 '现有笔记工具（Notion、Obsidian）学习成本高，缺乏智能整理能力。用户记了大量笔记但很少回顾，知识利用率低。',
 1, 20.00, 'standard', 'active', '2026-02-18 10:00:00+08');

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

(3, 1, 'D3', '产品需求文档',
 'Create a product requirements document with user stories, acceptance criteria, and feature prioritization.',
 'User story format, P0 features <=10, includes IA and core flow diagrams',
 'D1,D2', 120000, 'submitted', 4, '2026-03-01 09:00:00+08', '2026-03-05 17:00:00+08', NULL,
 'PRD with 8 P0 features, 12 P1 features, information architecture, and 4 core user flows.', NULL, 0.980000),

(4, 1, 'D4', '技术可行性评估',
 'Evaluate technical feasibility including technology stack recommendations and key risk points.',
 'Tech stack recommendations + architecture overview + risk points + feasibility conclusion',
 'D3', 80000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0);

-- ============================================================
-- Idea 2: 社区团购小程序 (active, D1 approved, D2 claimed)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(2,
 '社区团购小程序 — 邻里拼',
 '面向小区居民的社区团购平台。团长发起拼团，邻居一键参团，支持自提点取货。解决生鲜水果最后一公里配送成本高的问题。',
 '社区团长、小区居民、生鲜供应商',
 '现有社区团购平台（美团优选、多多买菜）佣金高、品类受限。小区自发团购用微信群管理混乱，对账困难。',
 2, 15.00, 'standard', 'active', '2026-02-25 14:00:00+08');

INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, claimed_by, claimed_at, submitted_at, approved_at, output_content, quality_score, cost_usd_accumulated) VALUES
(5, 2, 'D1', '竞品分析报告',
 'Analyze community group-buying platforms.',
 '>=3 direct competitors analyzed',
 '', 80000, 'approved', 4, '2026-02-26 10:00:00+08', '2026-02-28 15:00:00+08', '2026-03-01 09:00:00+08',
 'Analyzed 美团优选, 多多买菜, 淘菜菜, 兴盛优选, and 3 indie WeChat mini-programs.',
 1.00, 0.720000),

(6, 2, 'D2', '目标用户画像',
 'Define personas for group-buying participants.',
 '2-3 user personas with scenarios',
 '', 60000, 'claimed', 5, '2026-03-02 11:00:00+08', NULL, NULL,
 NULL, NULL, 0.080000),

(7, 2, 'D3', '产品需求文档',
 'Create PRD for the group-buying mini-program.',
 'User stories with acceptance criteria, P0 features <=10',
 'D1,D2', 120000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0),

(8, 2, 'D4', '技术可行性评估',
 'Evaluate tech stack for WeChat mini-program.',
 'Tech stack recommendations + risk points + feasibility conclusion',
 'D3', 80000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0);

-- ============================================================
-- Idea 3: 宠物健康管理 App (completed, all 4 approved)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(3,
 '宠物健康管理 App — PetCare+',
 '帮助宠物主人记录宠物的饮食、运动、疫苗、体检等健康数据，AI 分析异常指标并提醒就医。支持多宠物管理和兽医在线问诊。',
 '养猫养狗的年轻人、宠物店主、兽医',
 '宠物就医成本高且信息不对称。主人难以判断宠物是否生病，错过最佳治疗时机。现有宠物 App 功能碎片化，缺乏健康数据追踪。',
 3, 25.00, 'standard', 'completed', '2026-01-15 09:00:00+08');

INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, claimed_by, claimed_at, submitted_at, approved_at, quality_score, cost_usd_accumulated) VALUES
(9,  3, 'D1', '竞品分析报告', 'Analyze pet health management apps.', '>=3 competitors', '', 80000, 'approved', 2, '2026-01-16 10:00:00+08', '2026-01-18 16:00:00+08', '2026-01-19 09:00:00+08', 1.20, 0.900000),
(10, 3, 'D2', '目标用户画像', 'Define pet owner personas.', '2-3 personas', '', 60000, 'approved', 4, '2026-01-17 14:00:00+08', '2026-01-20 11:00:00+08', '2026-01-21 10:00:00+08', 1.50, 0.550000),
(11, 3, 'D3', '产品需求文档', 'Create PRD for pet health app.', 'User stories, P0 <=10, IA + flows', 'D1,D2', 120000, 'approved', 2, '2026-01-22 09:00:00+08', '2026-01-26 17:00:00+08', '2026-01-27 11:00:00+08', 1.20, 1.100000),
(12, 3, 'D4', '技术可行性评估', 'Evaluate tech stack for pet health platform.', 'Tech recommendations + risks', 'D3', 80000, 'approved', 3, '2026-01-28 10:00:00+08', '2026-02-01 15:00:00+08', '2026-02-02 09:00:00+08', 1.00, 0.650000);

-- ============================================================
-- Idea 4: 自由职业者接单平台 (active, all open)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(4,
 '自由职业者接单平台 — FreelanceHub',
 '连接自由职业者和企业的智能匹配平台。基于技能图谱和项目需求自动推荐合适的自由职业者，支持里程碑付款和在线协作。',
 '设计师、开发者、文案、翻译等自由职业者；中小企业、创业团队',
 '现有平台（猪八戒、Fiverr）匹配效率低、抽佣高。自由职业者获客难，企业找到合适人才耗时长。',
 4, 18.00, 'standard', 'active', '2026-03-08 11:00:00+08');

INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, cost_usd_accumulated) VALUES
(13, 4, 'D1', '竞品分析报告', 'Analyze freelance platforms.', '>=3 competitors', '', 80000, 'open', 0),
(14, 4, 'D2', '目标用户画像', 'Define freelancer and client personas.', '2-3 personas', '', 60000, 'open', 0),
(15, 4, 'D3', '产品需求文档', 'Create PRD for freelance matching platform.', 'User stories, P0 <=10', 'D1,D2', 120000, 'open', 0),
(16, 4, 'D4', '技术可行性评估', 'Evaluate tech stack for matching and payment.', 'Tech recommendations + risks', 'D3', 80000, 'open', 0);

-- ============================================================
-- Idea 5: 商标交易平台 (active, D1+D2 approved, D3 claimed)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(5,
 '商标交易平台 — 标易通',
 '面向中小企业和个体商户的商标交易撮合平台。整合闲置商标资源，通过 AI 估价、智能匹配和在线过户一站式服务，降低商标获取成本和时间。支持商标买卖、授权许可和质押融资。',
 '创业公司创始人、中小企业主、知识产权代理机构、个体工商户',
 '商标注册周期长（6-12个月）且通过率低（约40%）。现有商标交易平台信息不透明、估价混乱、过户流程复杂。中小企业急需商标时缺乏高效获取渠道。',
 1, 22.00, 'standard', 'active', '2026-03-05 09:00:00+08');

INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, claimed_by, claimed_at, submitted_at, approved_at, output_content, quality_score, cost_usd_accumulated) VALUES
(17, 5, 'D1', '竞品分析报告',
 'Research and analyze existing trademark trading platforms, IP marketplaces, and related services.',
 '>=3 direct competitors + >=2 indirect competitors analyzed',
 '', 80000, 'approved', 3, '2026-03-06 10:00:00+08', '2026-03-08 16:00:00+08', '2026-03-09 09:00:00+08',
 'Analyzed 中华商标超市, 八戒知产, 权大师, 标天下, and indirect competitors.',
 1.20, 0.880000),

(18, 5, 'D2', '目标用户画像',
 'Define target user personas for trademark buyers, sellers, and IP agencies.',
 '2-3 user personas with narrative scenarios',
 '', 60000, 'approved', 4, '2026-03-07 14:00:00+08', '2026-03-09 11:00:00+08', '2026-03-10 10:00:00+08',
 'Created 3 personas: startup founder, retiring merchant, IP agency.',
 1.50, 0.590000),

(19, 5, 'D3', '产品需求文档',
 'Create PRD including AI valuation, smart matching, and online transfer features.',
 'User stories, P0 features <=10, IA + core flows',
 'D1,D2', 120000, 'claimed', 2, '2026-03-10 09:00:00+08', NULL, NULL,
 NULL, NULL, 0.320000),

(20, 5, 'D4', '技术可行性评估',
 'Evaluate tech stack for trademark data integration and AI valuation model.',
 'Tech stack recommendations + risk points + feasibility conclusion',
 'D3', 80000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0);

SELECT setval('ideas_id_seq', 5);
SELECT setval('tasks_id_seq', 20);

-- ============================================================
-- Documents
-- ============================================================

-- Idea 1 documents
INSERT INTO documents (id, task_id, content, current_version) VALUES
(1, 1, E'# 竞品分析报告 — AI 笔记助手\n\n## 1. 直接竞品\n\n### 1.1 Notion AI\n- **定位**: All-in-one workspace with AI assistant\n- **优势**: 生态完善，模板丰富，协作能力强\n- **劣势**: 学习曲线陡，AI 功能偏通用，不专注笔记场景\n- **定价**: $10/月（Plus），AI 附加 $8/月\n\n### 1.2 Mem\n- **定位**: AI-first note-taking for professionals\n- **优势**: 自动关联、智能搜索、零组织成本\n- **劣势**: 生态封闭，不支持本地存储，团队功能弱\n- **定价**: $15/月\n\n### 1.3 Reflect\n- **定位**: 端到端加密的 AI 笔记\n- **优势**: 隐私保护强，AI 摘要质量高\n- **劣势**: 功能简单，无协作，市场份额小\n\n## 2. 间接竞品\n\n### 2.1 Obsidian\n- 本地优先、插件生态强，但 AI 能力需第三方插件\n\n### 2.2 Apple Notes\n- 零学习成本，但智能化程度低\n\n## 3. 差异化空间\n- **语音转笔记**: 会议/灵感场景，竞品普遍弱\n- **知识图谱可视化**: 超越简单双链\n- **中文语境优化**: 国际竞品中文支持差', 2),
(2, 2, E'# 目标用户画像\n\n## Persona 1: 独立开发者 Alex\n- **年龄**: 28岁\n- **痛点**: 学习新技术时笔记分散在多个平台\n- **场景**: 边写代码边记 debug 日志\n- **预算**: ¥30-50/月\n\n## Persona 2: 研究生 小林\n- **年龄**: 24岁\n- **痛点**: 论文阅读笔记多但缺乏系统整理\n- **场景**: 课堂录音转文字，文献笔记自动关联\n- **预算**: ¥15-25/月\n\n## Persona 3: 产品经理 Jenny\n- **年龄**: 32岁\n- **痛点**: 需求调研、用户访谈、竞品分析散落各处\n- **场景**: 会议纪要自动生成 action items\n- **预算**: ¥50-100/月', 2),
(3, 3, '', 1),
(4, 4, '', 1);

-- Idea 2 documents
INSERT INTO documents (id, task_id, content, current_version) VALUES
(5, 5, E'# 竞品分析 — 社区团购\n\n## 直接竞品\n1. **美团优选** — 流量大但佣金 15-20%\n2. **多多买菜** — 低价策略，依赖拼多多流量\n3. **淘菜菜** — 阿里系供应链\n4. **兴盛优选** — 下沉市场强\n\n## 差异化\n- 去中心化：团长自主定价\n- 本地化供应链对接\n- 低佣金（5%）吸引团长', 1),
(6, 6, '', 1),
(7, 7, '', 1),
(8, 8, '', 1);

-- Idea 3 documents (all approved)
INSERT INTO documents (id, task_id, content, current_version) VALUES
(9,  9,  '# 竞品分析 — 宠物健康管理', 1),
(10, 10, '# 目标用户画像', 1),
(11, 11, '# 产品需求文档', 1),
(12, 12, '# 技术可行性评估', 1);

-- Idea 4 documents (all empty)
INSERT INTO documents (id, task_id, content, current_version) VALUES
(13, 13, '', 1),
(14, 14, '', 1),
(15, 15, '', 1),
(16, 16, '', 1);

-- Idea 5 documents
INSERT INTO documents (id, task_id, content, current_version) VALUES
(17, 17, E'# 竞品分析报告 — 商标交易平台\n\n## 1. 直接竞品\n\n### 1.1 中华商标超市网\n- **定位**: 国内最大商标转让平台\n- **优势**: 商标库量大（200万+）\n- **劣势**: 界面老旧，估价不透明\n\n### 1.2 八戒知产\n- **定位**: 一站式知识产权服务\n- **优势**: 流量大，服务链完整\n- **劣势**: 商标交易非核心业务\n\n### 1.3 权大师\n- **定位**: 智能知识产权服务平台\n- **优势**: AI 辅助检索\n- **劣势**: 侧重注册，交易功能弱\n\n## 2. 差异化空间\n- AI 估价模型\n- 智能匹配\n- 过户全托管', 2),
(18, 18, E'# 目标用户画像 — 商标交易\n\n## Persona 1: 创业公司创始人 李明\n- 注册商标周期太长，直接购买已注册商标\n- 预算: 5000-30000元\n\n## Persona 2: 个体商户 王姐\n- 闲置商标想变现\n- 希望零成本挂牌\n\n## Persona 3: 知识产权代理人 张律师\n- 管理 500+ 客户商标，需批量工具\n- 愿付 ¥200-500/月 SaaS 费用', 2),
(19, 19, '', 1),
(20, 20, '', 1);

SELECT setval('documents_id_seq', 20);

-- ============================================================
-- Token usage logs
-- ============================================================
INSERT INTO token_usage_logs (user_id, task_id, model, tokens_in, tokens_out, cost_usd, timestamp) VALUES
-- Idea 1
(2, 1, 'claude-3.5-sonnet', 45000, 12000, 0.420000, '2026-02-20 14:00:00+08'),
(2, 1, 'claude-3.5-sonnet', 38000, 15000, 0.430000, '2026-02-20 17:00:00+08'),
(3, 2, 'claude-3.5-sonnet', 32000, 18000, 0.620000, '2026-02-22 10:00:00+08'),
(4, 3, 'claude-3.5-sonnet', 52000, 20000, 0.980000, '2026-03-05 15:00:00+08'),
-- Idea 2
(4, 5, 'claude-3.5-sonnet', 40000, 16000, 0.720000, '2026-02-28 14:00:00+08'),
(5, 6, 'claude-3.5-sonnet', 8000,  3000,  0.080000, '2026-03-03 11:00:00+08'),
-- Idea 3
(2, 9,  'claude-3.5-sonnet', 50000, 14000, 0.900000, '2026-01-18 15:00:00+08'),
(4, 10, 'claude-3.5-sonnet', 30000, 16000, 0.550000, '2026-01-20 10:00:00+08'),
(2, 11, 'claude-3.5-sonnet', 60000, 22000, 1.100000, '2026-01-26 16:00:00+08'),
(3, 12, 'claude-3.5-sonnet', 35000, 15000, 0.650000, '2026-02-01 14:00:00+08'),
-- Idea 5
(3, 17, 'claude-3.5-sonnet', 48000, 15000, 0.880000, '2026-03-08 15:00:00+08'),
(4, 18, 'claude-3.5-sonnet', 30000, 16000, 0.590000, '2026-03-09 10:00:00+08'),
(2, 19, 'claude-3.5-sonnet', 18000, 8000,  0.320000, '2026-03-11 14:00:00+08');

-- ============================================================
-- Contributions (for approved tasks)
-- ============================================================
INSERT INTO contributions (idea_id, task_id, user_id, cost_usd, quality_score, weighted_score, weight_percent) VALUES
-- Idea 1
(1, 1, 2, 0.850000, 1.20, 1.020000, 0),
(1, 2, 3, 0.620000, 1.50, 0.930000, 0),
-- Idea 2
(2, 5, 4, 0.720000, 1.00, 0.720000, 100.0000),
-- Idea 3
(3, 9,  2, 0.900000, 1.20, 1.080000, 28.0000),
(3, 10, 4, 0.550000, 1.50, 0.825000, 21.4000),
(3, 11, 2, 1.100000, 1.20, 1.320000, 34.2000),
(3, 12, 3, 0.650000, 1.00, 0.650000, 16.4000),
-- Idea 5
(5, 17, 3, 0.880000, 1.20, 1.056000, 0),
(5, 18, 4, 0.590000, 1.50, 0.885000, 0);

-- ============================================================
-- PRD for completed idea 3
-- ============================================================
INSERT INTO prds (id, idea_id, content, published_at, price_credits, read_count) VALUES
(1, 3,
 E'# PetCare+ — 宠物健康管理 App PRD\n\n## 产品概述\n面向年轻宠物主人的一站式健康管理工具\n\n## 核心功能\n1. 健康数据记录（体重、饮食、运动）\n2. 疫苗与体检提醒\n3. AI 异常指标预警\n4. 兽医在线问诊\n5. 多宠物管理\n\n## 商业模式\n- 基础版免费\n- Pro 版 ¥25/月（AI 分析 + 兽医问诊）',
 '2026-02-10 09:00:00+08', 500.0000, 3);

SELECT setval('prds_id_seq', 1);

-- ============================================================
-- Credit transactions
-- ============================================================
INSERT INTO credit_transactions (user_id, type, amount, reference_type, reference_id, description, created_at) VALUES
(2, 'earn_contribute', 1020.0000, 'contribution', 1, 'D1 竞品分析 — AI 笔记助手', '2026-02-21 10:00:00+08'),
(3, 'earn_contribute', 930.0000, 'contribution', 2, 'D2 用户画像 — AI 笔记助手', '2026-02-23 09:00:00+08'),
(4, 'earn_contribute', 720.0000, 'contribution', 3, 'D1 竞品分析 — 社区团购', '2026-03-01 09:00:00+08'),
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
