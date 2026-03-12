-- Seed data for Claway (4 tasks per idea: doc1-doc4)
-- 10 users, 8 ideas, diverse statuses, 2 published PRDs
BEGIN;

-- ============================================================
-- Users (10 users with diverse roles)
-- ============================================================
INSERT INTO users (id, openclaw_id, username, display_name, avatar_url, credits_balance, created_at) VALUES
(1,  'oc_u_a1c3chen',    'alicechen_dev',   'Alice Chen',     '', 2520.0000, '2026-01-10 08:00:00+08'),
(2,  'oc_u_b0bzhang',    'bobzhang42',      'Bob Zhang',      '', 4800.5000, '2026-01-12 10:30:00+08'),
(3,  'oc_u_car0lwang',   'carolwang_pm',    'Carol Wang',     '', 1860.0000, '2026-01-15 14:00:00+08'),
(4,  'oc_u_dav1dli',     'davidli_ux',      'David Li',       '',  950.0000, '2026-01-20 09:00:00+08'),
(5,  'oc_u_3v3liu',      'eveliu_eng',      'Eve Liu',        '',  320.0000, '2026-02-01 16:00:00+08'),
(6,  'oc_u_fr4nkxu',     'frankxu_ops',     'Frank Xu',       '', 1200.0000, '2026-02-05 11:00:00+08'),
(7,  'oc_u_gr4cezhao',   'gracezhao_ai',    'Grace Zhao',     '',  680.0000, '2026-02-10 09:30:00+08'),
(8,  'oc_u_h3nrysun',    'henrysun_data',   'Henry Sun',      '',  150.0000, '2026-02-15 14:00:00+08'),
(9,  'oc_u_1r1shuang',   'irishuang_mkt',   'Iris Huang',     '',    0.0000, '2026-02-20 10:00:00+08'),
(10, 'oc_u_j4ckwu',      'jackwu_arch',     'Jack Wu',        '', 3500.0000, '2026-03-01 08:00:00+08');

SELECT setval('users_id_seq', 10);

-- ============================================================
-- Idea 1: AI 笔记助手 (active, doc1+doc2 approved, doc3 in revision, doc4 open)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(1,
 'AI 笔记助手 — NoteFlow',
 '为独立开发者和学生设计的智能笔记工具。支持语音转文字、自动分类、跨笔记关联搜索和 AI 摘要生成。目标是让碎片化知识变成可检索、可复用的知识库。',
 '独立开发者、大学生、知识工作者',
 '现有笔记工具（Notion、Obsidian）学习成本高，缺乏智能整理能力。用户记了大量笔记但很少回顾，知识利用率低。',
 1, 20.00, 'standard', 'active', '2026-01-18 10:00:00+08');

INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, claimed_by, claimed_at, submitted_at, approved_at, output_content, quality_score, review_feedback, cost_usd_accumulated) VALUES
(1, 1, 'doc1', '竞品分析报告',
 'Research and analyze the competitive landscape for AI-powered note-taking tools.',
 '>=3 direct competitors + >=2 indirect competitors analyzed',
 '', 80000, 'approved', 2, '2026-01-19 09:00:00+08', '2026-01-20 18:00:00+08', '2026-01-21 10:00:00+08',
 'Completed competitive analysis covering Notion AI, Mem, Reflect, Obsidian, and Apple Notes.',
 1.20, NULL, 0.850000),

(2, 1, 'doc2', '目标用户画像',
 'Define target user personas with core pain points and usage scenarios.',
 '2-3 user personas with narrative scenarios',
 '', 60000, 'approved', 3, '2026-01-20 14:00:00+08', '2026-01-22 11:00:00+08', '2026-01-23 09:00:00+08',
 'Created 3 detailed personas: indie developer, graduate student, product manager.',
 1.50, NULL, 0.620000),

(3, 1, 'doc3', '产品需求文档',
 'Create a product requirements document with user stories, acceptance criteria, and feature prioritization.',
 'User story format, P0 features <=10, includes IA and core flow diagrams',
 'doc1,doc2', 120000, 'revision', 4, '2026-02-01 09:00:00+08', '2026-02-05 17:00:00+08', NULL,
 'PRD with 8 P0 features, 12 P1 features, information architecture, and 4 core user flows.',
 NULL, '用户故事的验收标准不够具体，请为每个 P0 功能补充可量化的验收标准。另外信息架构图缺少设置页面的子结构。', 0.980000),

(4, 1, 'doc4', '技术可行性评估',
 'Evaluate technical feasibility including technology stack recommendations and key risk points.',
 'Tech stack recommendations + architecture overview + risk points + feasibility conclusion',
 'doc3', 80000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0);

-- ============================================================
-- Idea 2: 社区团购小程序 (active, doc1 approved, doc2 submitted)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(2,
 '社区团购小程序 — 邻里拼',
 '面向小区居民的社区团购平台。团长发起拼团，邻居一键参团，支持自提点取货。解决生鲜水果最后一公里配送成本高的问题。',
 '社区团长、小区居民、生鲜供应商',
 '现有社区团购平台（美团优选、多多买菜）佣金高、品类受限。小区自发团购用微信群管理混乱，对账困难。',
 2, 15.00, 'standard', 'active', '2026-01-25 14:00:00+08');

INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, claimed_by, claimed_at, submitted_at, approved_at, output_content, quality_score, cost_usd_accumulated) VALUES
(5, 2, 'doc1', '竞品分析报告',
 'Analyze community group-buying platforms.',
 '>=3 direct competitors analyzed',
 '', 80000, 'approved', 4, '2026-01-26 10:00:00+08', '2026-01-28 15:00:00+08', '2026-02-01 09:00:00+08',
 'Analyzed 美团优选, 多多买菜, 淘菜菜, 兴盛优选, and 3 indie WeChat mini-programs.',
 1.00, 0.720000),

(6, 2, 'doc2', '目标用户画像',
 'Define personas for group-buying participants.',
 '2-3 user personas with scenarios',
 '', 60000, 'submitted', 5, '2026-02-02 11:00:00+08', '2026-02-06 16:00:00+08', NULL,
 'Created personas for community leader, busy parent, and fresh produce supplier.', NULL, 0.480000),

(7, 2, 'doc3', '产品需求文档',
 'Create PRD for the group-buying mini-program.',
 'User stories with acceptance criteria, P0 features <=10',
 'doc1,doc2', 120000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0),

(8, 2, 'doc4', '技术可行性评估',
 'Evaluate tech stack for WeChat mini-program.',
 'Tech stack recommendations + risk points + feasibility conclusion',
 'doc3', 80000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0);

-- ============================================================
-- Idea 3: 宠物健康管理 App (completed, all 4 approved, PRD published)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(3,
 '宠物健康管理 App — PetCare+',
 '帮助宠物主人记录宠物的饮食、运动、疫苗、体检等健康数据，AI 分析异常指标并提醒就医。支持多宠物管理和兽医在线问诊。',
 '养猫养狗的年轻人、宠物店主、兽医',
 '宠物就医成本高且信息不对称。主人难以判断宠物是否生病，错过最佳治疗时机。现有宠物 App 功能碎片化，缺乏健康数据追踪。',
 3, 25.00, 'standard', 'completed', '2026-01-05 09:00:00+08');

INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, claimed_by, claimed_at, submitted_at, approved_at, quality_score, cost_usd_accumulated) VALUES
(9,  3, 'doc1', '竞品分析报告', 'Analyze pet health management apps.', '>=3 competitors', '', 80000, 'approved', 2, '2026-01-06 10:00:00+08', '2026-01-08 16:00:00+08', '2026-01-09 09:00:00+08', 1.20, 0.900000),
(10, 3, 'doc2', '目标用户画像', 'Define pet owner personas.', '2-3 personas', '', 60000, 'approved', 6, '2026-01-07 14:00:00+08', '2026-01-10 11:00:00+08', '2026-01-11 10:00:00+08', 1.50, 0.550000),
(11, 3, 'doc3', '产品需求文档', 'Create PRD for pet health app.', 'User stories, P0 <=10, IA + flows', 'doc1,doc2', 120000, 'approved', 2, '2026-01-12 09:00:00+08', '2026-01-16 17:00:00+08', '2026-01-17 11:00:00+08', 1.20, 1.100000),
(12, 3, 'doc4', '技术可行性评估', 'Evaluate tech stack for pet health platform.', 'Tech recommendations + risks', 'doc3', 80000, 'approved', 7, '2026-01-18 10:00:00+08', '2026-01-22 15:00:00+08', '2026-01-23 09:00:00+08', 1.00, 0.650000);

-- ============================================================
-- Idea 4: 自由职业者接单平台 (active, all open — fresh idea)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(4,
 '自由职业者接单平台 — FreelanceHub',
 '连接自由职业者和企业的智能匹配平台。基于技能图谱和项目需求自动推荐合适的自由职业者，支持里程碑付款和在线协作。',
 '设计师、开发者、文案、翻译等自由职业者；中小企业、创业团队',
 '现有平台（猪八戒、Fiverr）匹配效率低、抽佣高。自由职业者获客难，企业找到合适人才耗时长。',
 4, 18.00, 'standard', 'active', '2026-03-08 11:00:00+08');

INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, cost_usd_accumulated) VALUES
(13, 4, 'doc1', '竞品分析报告', 'Analyze freelance platforms.', '>=3 competitors', '', 80000, 'open', 0),
(14, 4, 'doc2', '目标用户画像', 'Define freelancer and client personas.', '2-3 personas', '', 60000, 'open', 0),
(15, 4, 'doc3', '产品需求文档', 'Create PRD for freelance matching platform.', 'User stories, P0 <=10', 'doc1,doc2', 120000, 'open', 0),
(16, 4, 'doc4', '技术可行性评估', 'Evaluate tech stack for matching and payment.', 'Tech recommendations + risks', 'doc3', 80000, 'open', 0);

-- ============================================================
-- Idea 5: 商标交易平台 (active, doc1+doc2 approved, doc3 claimed)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(5,
 '商标交易平台 — 标易通',
 '面向中小企业和个体商户的商标交易撮合平台。整合闲置商标资源，通过 AI 估价、智能匹配和在线过户一站式服务，降低商标获取成本和时间。',
 '创业公司创始人、中小企业主、知识产权代理机构、个体工商户',
 '商标注册周期长（6-12个月）且通过率低（约40%）。现有商标交易平台信息不透明、估价混乱、过户流程复杂。',
 1, 22.00, 'standard', 'active', '2026-02-05 09:00:00+08');

INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, claimed_by, claimed_at, submitted_at, approved_at, output_content, quality_score, cost_usd_accumulated) VALUES
(17, 5, 'doc1', '竞品分析报告',
 'Research and analyze existing trademark trading platforms.',
 '>=3 direct competitors + >=2 indirect competitors analyzed',
 '', 80000, 'approved', 3, '2026-02-06 10:00:00+08', '2026-02-08 16:00:00+08', '2026-02-09 09:00:00+08',
 'Analyzed 中华商标超市, 八戒知产, 权大师, 标天下, and indirect competitors.',
 1.20, 0.880000),

(18, 5, 'doc2', '目标用户画像',
 'Define target user personas for trademark buyers, sellers, and IP agencies.',
 '2-3 user personas with narrative scenarios',
 '', 60000, 'approved', 8, '2026-02-07 14:00:00+08', '2026-02-09 11:00:00+08', '2026-02-10 10:00:00+08',
 'Created 3 personas: startup founder, retiring merchant, IP agency.',
 1.50, 0.590000),

(19, 5, 'doc3', '产品需求文档',
 'Create PRD including AI valuation, smart matching, and online transfer features.',
 'User stories, P0 features <=10, IA + core flows',
 'doc1,doc2', 120000, 'claimed', 10, '2026-02-11 09:00:00+08', NULL, NULL,
 NULL, NULL, 0.320000),

(20, 5, 'doc4', '技术可行性评估',
 'Evaluate tech stack for trademark data integration and AI valuation model.',
 'Tech stack recommendations + risk points + feasibility conclusion',
 'doc3', 80000, 'open', NULL, NULL, NULL, NULL, NULL, NULL, 0);

-- ============================================================
-- Idea 6: 在线教育直播平台 (completed, all 4 approved, PRD published)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(6,
 '在线教育直播平台 — LiveClass',
 '面向 K12 和成人教育的互动直播教学平台。支持白板协作、实时答题、AI 助教和课后回放。帮助教育机构低成本转型线上。',
 '教育机构、独立讲师、K12 学生家长、职业培训学员',
 '现有直播平台（腾讯课堂、ClassIn）价格高、互动功能弱。小型教育机构难以承担高额平台费用，学生互动体验差。',
 6, 20.00, 'standard', 'completed', '2025-12-15 09:00:00+08');

INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, claimed_by, claimed_at, submitted_at, approved_at, quality_score, cost_usd_accumulated) VALUES
(21, 6, 'doc1', '竞品分析报告', 'Analyze online education live streaming platforms.', '>=3 competitors', '', 80000, 'approved', 4, '2025-12-16 10:00:00+08', '2025-12-18 16:00:00+08', '2025-12-19 09:00:00+08', 1.00, 0.780000),
(22, 6, 'doc2', '目标用户画像', 'Define personas for teachers, students, and institutions.', '2-3 personas', '', 60000, 'approved', 7, '2025-12-17 14:00:00+08', '2025-12-20 11:00:00+08', '2025-12-21 10:00:00+08', 1.20, 0.520000),
(23, 6, 'doc3', '产品需求文档', 'Create PRD for live education platform.', 'User stories, P0 <=10, IA + flows', 'doc1,doc2', 120000, 'approved', 10, '2025-12-22 09:00:00+08', '2025-12-26 17:00:00+08', '2025-12-27 11:00:00+08', 1.50, 1.250000),
(24, 6, 'doc4', '技术可行性评估', 'Evaluate tech stack for live streaming and real-time interaction.', 'Tech recommendations + risks', 'doc3', 80000, 'approved', 2, '2025-12-28 10:00:00+08', '2026-01-02 15:00:00+08', '2026-01-03 09:00:00+08', 1.20, 0.750000);

-- ============================================================
-- Idea 7: 健身饮食规划 App (active, doc1 claimed, doc2 open)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(7,
 '健身饮食规划 App — FitMeal',
 '基于用户身体数据和健身目标，AI 自动生成个性化饮食计划和食谱推荐。支持拍照记录饮食、营养成分分析和超市购物清单生成。',
 '健身爱好者、减脂增肌人群、健身教练',
 '现有饮食管理工具（薄荷健康、MyFitnessPal）手动记录繁琐，缺乏个性化推荐。用户难以坚持饮食计划，健身效果打折扣。',
 5, 15.00, 'standard', 'active', '2026-03-01 09:00:00+08');

INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, claimed_by, claimed_at, cost_usd_accumulated) VALUES
(25, 7, 'doc1', '竞品分析报告', 'Analyze fitness and diet planning apps.', '>=3 competitors', '', 80000, 'claimed', 6, '2026-03-02 10:00:00+08', 0.150000),
(26, 7, 'doc2', '目标用户画像', 'Define fitness enthusiast personas.', '2-3 personas', '', 60000, 'open', NULL, NULL, 0),
(27, 7, 'doc3', '产品需求文档', 'Create PRD for diet planning features.', 'User stories, P0 <=10', 'doc1,doc2', 120000, 'open', NULL, NULL, 0),
(28, 7, 'doc4', '技术可行性评估', 'Evaluate tech stack for AI diet recommendations.', 'Tech recommendations + risks', 'doc3', 80000, 'open', NULL, NULL, 0);

-- ============================================================
-- Idea 8: 二手车估价平台 (cancelled)
-- ============================================================
INSERT INTO ideas (id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at) VALUES
(8,
 '二手车估价平台 — 车价通',
 '基于大数据和 AI 的二手车精准估价工具。整合车辆维修记录、事故历史、市场行情等多维数据，为买卖双方提供透明公正的估价服务。',
 '二手车买家、卖家、车商、保险公司',
 '二手车市场信息严重不对称，估价主观性强。买家担心买贵，卖家担心卖亏。现有估价工具精度低，缺乏权威性。',
 8, 20.00, 'standard', 'cancelled', '2026-02-20 09:00:00+08');

INSERT INTO tasks (id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status, cost_usd_accumulated) VALUES
(29, 8, 'doc1', '竞品分析报告', 'Analyze used car valuation platforms.', '>=3 competitors', '', 80000, 'open', 0),
(30, 8, 'doc2', '目标用户画像', 'Define buyer and seller personas.', '2-3 personas', '', 60000, 'open', 0),
(31, 8, 'doc3', '产品需求文档', 'Create PRD for car valuation features.', 'User stories, P0 <=10', 'doc1,doc2', 120000, 'open', 0),
(32, 8, 'doc4', '技术可行性评估', 'Evaluate tech stack for AI valuation model.', 'Tech recommendations + risks', 'doc3', 80000, 'open', 0);

SELECT setval('ideas_id_seq', 8);
SELECT setval('tasks_id_seq', 32);

-- ============================================================
-- Documents
-- ============================================================

-- Idea 1 documents
INSERT INTO documents (id, task_id, content, current_version) VALUES
(1, 1, E'# 竞品分析报告 — AI 笔记助手\n\n## 执行摘要\n\nAI 笔记工具市场正在快速增长，但尚未出现兼具「低学习成本」和「深度智能化」的产品。现有竞品要么功能全面但复杂（Notion AI），要么智能但封闭（Mem），留下了明确的差异化空间。\n\n## 1. 直接竞品\n\n### 1.1 Notion AI\n- **定位**: All-in-one workspace with AI assistant\n- **目标用户**: 团队协作、项目管理\n- **优势**: 生态完善，模板丰富（10,000+），协作能力强，品牌认知度高\n- **劣势**: 学习曲线陡（新用户平均需 2 周上手），AI 功能偏通用，笔记搜索不够智能\n- **定价**: Free / Plus $10/月 / AI 附加 $8/月\n- **来源**: https://notion.so\n\n### 1.2 Mem\n- **定位**: AI-first note-taking for professionals\n- **目标用户**: 知识工作者、研究员\n- **优势**: 自动关联、智能搜索、零组织成本，「just dump it」理念\n- **劣势**: 生态封闭，不支持本地存储，团队功能弱，导出选项少\n- **定价**: $15/月（Personal），$20/月（Teams）\n- **来源**: https://mem.ai\n\n### 1.3 Reflect\n- **定位**: 端到端加密的 AI 笔记\n- **目标用户**: 隐私敏感用户\n- **优势**: E2E 加密，AI 摘要质量高，日历集成\n- **劣势**: 功能简单，无协作，市场份额小（<50K 用户）\n- **定价**: $10/月\n- **来源**: https://reflect.app\n\n## 2. 间接竞品\n\n### 2.1 Obsidian\n- **定位**: 本地优先的知识管理工具\n- **优势**: 插件生态丰富（1,500+），双链笔记先驱，完全免费（个人）\n- **劣势**: AI 能力需第三方插件，学习成本不低，移动端体验一般\n- **来源**: https://obsidian.md\n\n### 2.2 Apple Notes\n- **定位**: 零学习成本的系统级笔记\n- **优势**: 开箱即用，iCloud 同步，手写支持\n- **劣势**: 智能化程度极低，无 AI 功能，跨平台差\n\n## 3. 对比表\n\n| 维度 | Notion AI | Mem | Reflect | Obsidian | 我们的机会 |\n|------|-----------|-----|---------|----------|----------|\n| 学习成本 | 高 | 中 | 低 | 中高 | 极低 |\n| AI 智能度 | 中 | 高 | 中 | 低 | 高 |\n| 语音支持 | 无 | 无 | 无 | 插件 | 核心功能 |\n| 中文优化 | 弱 | 无 | 无 | 弱 | 核心优势 |\n| 价格 | $18/月 | $15/月 | $10/月 | 免费 | ¥25-35/月 |\n\n## 4. 差异化空间\n\n1. **语音转笔记**: 会议/灵感场景，竞品普遍缺失\n2. **知识图谱可视化**: 超越简单双链，展示知识关联网络\n3. **中文语境优化**: 分词、语义理解、行业术语\n4. **极低学习成本**: 不需要教程就能上手\n\n## 结论\n\n市场存在明确机会——为中文用户打造一款「开箱即用 + 深度智能」的笔记工具。关键差异化在语音入口和中文语义理解。', 3),

(2, 2, E'# 目标用户画像 — AI 笔记助手\n\n## 概述\n\n基于竞品分析和目标市场调研，定义三类核心用户画像。主要目标用户是「有大量碎片化知识需要整理，但缺乏时间和工具」的知识工作者。\n\n## Persona 1: 独立开发者 Alex\n\n### 人口统计\n- **角色**: 全栈独立开发者\n- **年龄**: 25-32 岁\n- **技术素养**: 高\n- **收入**: ¥15,000-30,000/月\n\n### 目标\n- 高效记录技术学习笔记和 debug 经验\n- 跨项目复用知识片段\n- 减少重复踩坑\n\n### 痛点\n- 笔记分散在 Notion、VS Code snippets、浏览器书签多处\n- 花大量时间学习工具本身而非专注内容\n- 搜索「那个之前解决过的 bug」总是找不到\n\n### 当前方案\n- Obsidian（主力） + 微信收藏（临时） + GitHub Issues（项目相关）\n\n### 场景 A: 学习新框架\nAlex 在 YouTube 看 React Server Components 教程，边看边用语音记录关键点。AI 自动将语音转为结构化笔记，并关联他之前写的 Next.js 笔记。三天后写代码时搜索「RSC data fetching」，直接找到当时的学习笔记。\n\n### 场景 B: Debug 日志\n凌晨 2 点定位到一个诡异的 CORS 问题。Alex 用语音快速描述问题和解决方案。两个月后遇到类似问题，AI 自动推荐「你之前解决过类似的 CORS 问题」。\n\n### 当前方案局限\n- Obsidian 需要手动整理和打标签，学习新插件成本高\n- 没有语音入口，灵感稍纵即逝\n- 跨工具搜索几乎不可能\n\n## Persona 2: 研究生 小林\n\n### 人口统计\n- **角色**: 计算机科学硕士研究生\n- **年龄**: 22-26 岁\n- **技术素养**: 中高\n- **预算**: ¥15-25/月\n\n### 目标\n- 高效整理论文阅读笔记\n- 课堂录音自动转文字\n- 为毕业论文积累素材\n\n### 痛点\n- 论文阅读笔记量大但缺乏系统整理\n- 课堂上来不及记笔记，课后又忘了重点\n- 临近毕业论文时发现之前的阅读笔记散落各处\n\n### 场景 A: 论文精读\n小林在 iPad 上阅读一篇关于 Transformer 优化的论文，用语音记录自己的理解和疑问。AI 自动提取关键发现，并关联他之前读过的 5 篇相关论文的笔记，生成一份「Transformer 优化技术演进」的知识图谱。\n\n### 场景 B: 课堂录音\n导师组会上，小林开启录音。会后 AI 自动生成会议纪要，提取 action items，并标记与自己课题相关的讨论点。\n\n### 当前方案局限\n- Zotero 管理论文但笔记功能弱\n- 手动整理耗时，经常「存了就忘」\n- 缺乏跨论文的知识关联能力\n\n## Persona 3: 产品经理 Jenny\n\n### 人口统计\n- **角色**: 中型 SaaS 公司产品经理\n- **年龄**: 28-35 岁\n- **技术素养**: 中\n- **预算**: ¥50-100/月（可报销）\n\n### 目标\n- 高效整理需求调研、用户访谈、竞品分析笔记\n- 会议纪要自动生成 action items\n- 快速回顾历史决策背景\n\n### 痛点\n- 需求调研、用户访谈、竞品分析散落在飞书文档、Notion、微信消息里\n- 开会记笔记影响参与讨论\n- 半年后回顾「为什么当时做了这个决定」很难找到原始讨论记录\n\n### 场景 A: 用户访谈\nJenny 和客户进行 60 分钟视频访谈。全程录音，AI 自动提取用户反馈要点、情感倾向和产品改进建议。访谈结束后 5 分钟就能分享整理好的访谈摘要给团队。\n\n### 场景 B: 需求评审\n在需求评审会上，Jenny 发现设计师对某个交互方案有疑问。她搜索「登录页改版讨论」，AI 找到三个月前的会议记录和当时的决策理由，当场解答。\n\n### 当前方案局限\n- 飞书会议纪要需要手动整理\n- 跨工具搜索不可能（飞书 + Notion + 微信）\n- 缺乏结构化的知识积累\n\n## 跨画像分析\n\n### 共同痛点\n1. 知识碎片化——信息散落在多个工具中\n2. 手动整理成本高——缺乏自动化\n3. 搜索能力弱——找不到「曾经记过的东西」\n\n### 差异化需求\n- Alex: 代码片段支持、技术搜索\n- 小林: 学术论文集成、知识图谱\n- Jenny: 团队共享、会议纪要\n\n## 结论\n\n首选主打 **Alex（独立开发者）** 作为初始用户群体。原因：\n1. 对新工具接受度最高\n2. 能产生口碑传播\n3. 付费意愿明确\n4. 使用场景清晰（学习 + debug）', 3),

(3, 3, E'# 产品需求文档 — AI 笔记助手（修订中）\n\n## 产品概述\n基于竞品分析（doc1）发现的市场空白和用户画像（doc2）定义的核心需求...\n\n## P0 功能\n1. 语音输入转结构化笔记\n2. AI 智能分类和标签\n3. 全文搜索 + 语义搜索\n4. 笔记间自动关联\n5. 移动端 + Web 端同步\n6. Markdown 编辑器\n7. 导入/导出（Markdown, JSON）\n8. 用户账户和登录\n\n（待补充验收标准和信息架构...）', 2),

(4, 4, '', 1);

-- Idea 2 documents
INSERT INTO documents (id, task_id, content, current_version) VALUES
(5, 5, E'# 竞品分析 — 社区团购\n\n## 执行摘要\n\n社区团购赛道经历了 2020-2022 年的「百团大战」后，巨头已建立流量壁垒，但在社区自治和低佣金领域仍有空间。\n\n## 直接竞品\n\n### 1. 美团优选\n- **模式**: 平台型，次日达\n- **佣金**: 15-20%\n- **优势**: 流量大，配送网络完善\n- **劣势**: 佣金高，团长收入下降，品类以标品为主\n\n### 2. 多多买菜\n- **模式**: 平台型，低价策略\n- **佣金**: 12-18%\n- **优势**: 价格极低，依赖拼多多流量\n- **劣势**: 品质不稳定，退货率高\n\n### 3. 淘菜菜\n- **模式**: 平台型，阿里系供应链\n- **佣金**: 15%\n- **优势**: 供应链强\n- **劣势**: 用户粘性差，2023 年已缩减城市\n\n### 4. 兴盛优选\n- **模式**: 下沉市场\n- **优势**: 三四线城市渗透率高\n- **劣势**: 技术能力弱\n\n## 差异化方向\n- 去中心化：团长自主定价和选品\n- 极低佣金（5%）吸引优质团长\n- 本地化供应链直连\n- 微信生态内完成全流程', 2),

(6, 6, E'# 目标用户画像 — 社区团购\n\n## Persona 1: 社区团长 张姐\n- **年龄**: 35-45 岁，全职妈妈\n- **动机**: 利用空闲时间赚取佣金\n- **痛点**: 用微信群管理拼团，接龙统计、收款对账全靠人工\n- **日常**: 每天在 3 个微信群发布 5-8 款商品，手动统计订单 2 小时\n\n## Persona 2: 上班族居民 李先生\n- **年龄**: 28-35 岁，互联网从业者\n- **动机**: 买到新鲜便宜的水果和蔬菜\n- **痛点**: 下班太晚超市关门，外卖配送费高\n- **日常**: 通勤路上浏览商品，下班回家路过自提点取货\n\n## Persona 3: 生鲜供应商 王总\n- **年龄**: 40-50 岁，水果批发商\n- **动机**: 拓展零售渠道，减少损耗\n- **痛点**: 大平台压价严重，账期长（45-60天）', 2),

(7, 7, '', 1),
(8, 8, '', 1);

-- Idea 3 documents (all approved, with rich content)
INSERT INTO documents (id, task_id, content, current_version) VALUES
(9,  9,  E'# 竞品分析 — 宠物健康管理\n\n## 直接竞品\n1. **小壹管家** — 国内宠物综合管理，功能全但 UI 老旧\n2. **PetDesk** — 美国市场，兽医预约为主\n3. **11pets** — 欧洲市场，健康记录为核心\n\n## 间接竞品\n1. **微信小程序「萌爪医生」** — 在线问诊\n2. **京东宠物** — 电商+服务\n\n## 差异化\n- AI 异常检测（体重、饮食模式变化预警）\n- 多宠物统一管理\n- 中国市场本地化（宠物医院数据接入）', 1),
(10, 10, E'# 用户画像 — 宠物健康管理\n\n## Persona 1: 新手铲屎官 小王\n- 养了第一只猫，焦虑且缺乏经验\n- 经常因为猫咪「不吃饭」「拉肚子」恐慌就医\n- 需要一个「宠物育儿手册」\n\n## Persona 2: 多宠家庭 陈姐\n- 2 猫 1 狗，疫苗时间总搞混\n- 需要统一管理多只宠物的健康档案\n\n## Persona 3: 宠物店主 赵老板\n- 管理 30+ 寄养宠物\n- 需要批量健康记录和提醒', 1),
(11, 11, E'# 产品需求文档 — PetCare+\n\n## P0 功能\n1. 宠物档案管理（多宠物）\n2. 健康数据记录（体重、饮食、运动）\n3. 疫苗/体检日历提醒\n4. AI 异常指标预警\n5. 兽医问诊入口\n\n## 核心用户流程\n1. 注册 → 添加宠物 → 填写基本信息\n2. 日常记录 → 体重/饮食/运动\n3. AI 分析 → 异常提醒 → 就医建议', 1),
(12, 12, E'# 技术可行性 — PetCare+\n\n## 技术栈推荐\n- **前端**: React Native（跨平台）\n- **后端**: Go + PostgreSQL\n- **AI**: 基于历史数据的异常检测模型\n\n## 关键风险\n1. 宠物医院数据接入——需要逐家对接\n2. AI 模型准确率——需要大量标注数据\n\n## 结论: 可行\nMVP 可在 3 个月内完成，AI 功能可后续迭代。', 1);

-- Idea 4 documents (all empty)
INSERT INTO documents (id, task_id, content, current_version) VALUES
(13, 13, '', 1),
(14, 14, '', 1),
(15, 15, '', 1),
(16, 16, '', 1);

-- Idea 5 documents
INSERT INTO documents (id, task_id, content, current_version) VALUES
(17, 17, E'# 竞品分析报告 — 商标交易平台\n\n## 1. 直接竞品\n\n### 1.1 中华商标超市网\n- **定位**: 国内最大商标转让平台\n- **商标库**: 200万+\n- **优势**: 数据量大，历史悠久\n- **劣势**: 界面老旧（2010 年代设计），估价不透明，搜索体验差\n\n### 1.2 八戒知产\n- **定位**: 一站式知识产权服务（猪八戒旗下）\n- **优势**: 流量大，服务链完整（注册+转让+维权）\n- **劣势**: 商标交易非核心业务，匹配效率低\n\n### 1.3 权大师\n- **定位**: 智能知识产权服务平台\n- **优势**: AI 辅助商标检索，近似度分析\n- **劣势**: 侧重注册，交易撮合功能弱\n\n### 1.4 标天下\n- **定位**: 商标交易+品牌服务\n- **优势**: 地推团队强\n- **劣势**: 线上体验差，佣金 15-20%\n\n## 2. 间接竞品\n- 阿里知识产权交易平台\n- 中国版权保护中心\n\n## 3. 差异化空间\n- AI 估价模型（基于类别、有效期、品牌知名度）\n- 智能匹配（需求方画像 × 商标特征）\n- 过户全托管（降低交易风险）\n- 低佣金（5-8% vs 行业 15-20%）', 2),

(18, 18, E'# 目标用户画像 — 商标交易\n\n## Persona 1: 创业公司创始人 李明\n- **年龄**: 30 岁，互联网创业者\n- **需求**: 公司急需上线品牌，注册商标周期太长（6-12 月）\n- **行为**: 愿意直接购买已注册商标\n- **预算**: 5,000-30,000 元\n- **痛点**: 现有平台搜索体验差，估价不透明，担心被坑\n\n## Persona 2: 个体商户 王姐\n- **年龄**: 45 岁，服装店老板\n- **需求**: 关了实体店，闲置商标想变现\n- **行为**: 希望零成本挂牌，有人买了再付佣金\n- **痛点**: 不懂线上操作，担心被骗\n\n## Persona 3: 知识产权代理人 张律师\n- **年龄**: 35 岁，IP 律所合伙人\n- **需求**: 管理 500+ 客户商标，需批量操作工具\n- **行为**: 愿付 ¥200-500/月 SaaS 费用\n- **痛点**: 现有工具无法批量管理，每次过户都要手动填表', 2),

(19, 19, '', 1),
(20, 20, '', 1);

-- Idea 6 documents (all approved)
INSERT INTO documents (id, task_id, content, current_version) VALUES
(21, 21, E'# 竞品分析 — 在线教育直播\n\n## 直接竞品\n1. **腾讯课堂** — 流量大但互动弱\n2. **ClassIn** — 互动强但价格高（¥5000+/年）\n3. **钉钉课堂** — 免费但功能简陋\n\n## 差异化\n- AI 助教（自动答疑、课堂总结）\n- 白板实时协作\n- 价格定位 ¥200-500/月', 1),
(22, 22, E'# 用户画像 — 在线教育直播\n\n## Persona 1: K12 数学老师 刘老师\n- 30 岁，公立学校教师，兼职线上辅导\n- 需要白板 + 录制回放功能\n\n## Persona 2: 职业培训讲师 陈sir\n- 40 岁，独立 Python 讲师\n- 需要代码实时演示 + 学生互动', 1),
(23, 23, E'# PRD — LiveClass\n\n## P0 功能\n1. 直播间创建和管理\n2. 白板实时协作\n3. 实时弹幕和答题\n4. 课后回放 + AI 摘要\n5. 学生签到和考勤\n\n## 核心流程\n1. 讲师创建课程 → 设置时间 → 分享链接\n2. 学生点击链接 → 进入直播间 → 互动\n3. 课后 → AI 生成课堂总结 → 学生回看', 1),
(24, 24, E'# 技术可行性 — LiveClass\n\n## 技术栈\n- **前端**: Next.js + WebRTC\n- **后端**: Go + Redis + PostgreSQL\n- **直播**: 声网 Agora SDK\n- **AI**: Claude API for 课堂总结\n\n## 关键风险\n1. WebRTC 兼容性（特别是移动端浏览器）\n2. 并发直播间性能（目标: 同时 100 间）\n\n## 结论: 可行\n借助声网 SDK 可大幅降低直播技术门槛。', 1);

-- Idea 7 documents
INSERT INTO documents (id, task_id, content, current_version) VALUES
(25, 25, '', 1),
(26, 26, '', 1),
(27, 27, '', 1),
(28, 28, '', 1);

-- Idea 8 documents
INSERT INTO documents (id, task_id, content, current_version) VALUES
(29, 29, '', 1),
(30, 30, '', 1),
(31, 31, '', 1),
(32, 32, '', 1);

SELECT setval('documents_id_seq', 32);

-- ============================================================
-- Document versions (for documents that have been edited)
-- ============================================================
INSERT INTO document_versions (document_id, version, content, diff_from_previous, created_at, created_by) VALUES
-- Idea 1, doc1 had 3 versions
(1, 1, '', NULL, '2026-01-19 09:00:00+08', 2),
(1, 2, 'Draft competitive analysis...', 'Initial draft', '2026-01-20 14:00:00+08', 2),
(1, 3, '# 竞品分析报告 — AI 笔记助手...', 'Final version', '2026-01-20 18:00:00+08', 2),
-- Idea 1, doc2 had 3 versions
(2, 1, '', NULL, '2026-01-20 14:00:00+08', 3),
(2, 2, 'Draft personas...', 'Initial draft', '2026-01-21 16:00:00+08', 3),
(2, 3, '# 目标用户画像...', 'Final version', '2026-01-22 11:00:00+08', 3),
-- Idea 1, doc3 had 2 versions (still in revision)
(3, 1, '', NULL, '2026-02-01 09:00:00+08', 4),
(3, 2, '# 产品需求文档...', 'First draft', '2026-02-05 17:00:00+08', 4),
-- Idea 2, doc2 had 2 versions
(6, 1, '', NULL, '2026-02-02 11:00:00+08', 5),
(6, 2, '# 目标用户画像 — 社区团购...', 'Completed draft', '2026-02-06 16:00:00+08', 5);

-- ============================================================
-- Token usage logs (realistic LLM usage records)
-- ============================================================
INSERT INTO token_usage_logs (user_id, task_id, model, tokens_in, tokens_out, cost_usd, timestamp) VALUES
-- Idea 1: AI 笔记助手
(2, 1, 'claude-sonnet-4-20250514', 45000, 12000, 0.420000, '2026-01-20 14:00:00+08'),
(2, 1, 'claude-sonnet-4-20250514', 38000, 15000, 0.430000, '2026-01-20 17:00:00+08'),
(3, 2, 'claude-sonnet-4-20250514', 32000, 18000, 0.620000, '2026-01-22 10:00:00+08'),
(4, 3, 'claude-sonnet-4-20250514', 52000, 20000, 0.980000, '2026-02-05 15:00:00+08'),

-- Idea 2: 社区团购
(4, 5, 'claude-sonnet-4-20250514', 40000, 16000, 0.720000, '2026-01-28 14:00:00+08'),
(5, 6, 'claude-sonnet-4-20250514', 28000, 12000, 0.400000, '2026-02-05 11:00:00+08'),
(5, 6, 'claude-haiku-4-5-20251001',  8000,  3000, 0.080000, '2026-02-06 16:00:00+08'),

-- Idea 3: 宠物健康管理
(2, 9,  'claude-sonnet-4-20250514', 50000, 14000, 0.900000, '2026-01-08 15:00:00+08'),
(6, 10, 'claude-sonnet-4-20250514', 30000, 16000, 0.550000, '2026-01-10 10:00:00+08'),
(2, 11, 'claude-sonnet-4-20250514', 60000, 22000, 1.100000, '2026-01-16 16:00:00+08'),
(7, 12, 'claude-sonnet-4-20250514', 35000, 15000, 0.650000, '2026-01-22 14:00:00+08'),

-- Idea 5: 商标交易
(3, 17, 'claude-sonnet-4-20250514', 48000, 15000, 0.880000, '2026-02-08 15:00:00+08'),
(8, 18, 'claude-sonnet-4-20250514', 30000, 16000, 0.590000, '2026-02-09 10:00:00+08'),
(10, 19, 'claude-sonnet-4-20250514', 18000, 8000, 0.320000, '2026-02-12 14:00:00+08'),

-- Idea 6: 在线教育直播
(4, 21,  'claude-sonnet-4-20250514', 42000, 14000, 0.780000, '2025-12-18 15:00:00+08'),
(7, 22,  'claude-sonnet-4-20250514', 28000, 12000, 0.520000, '2025-12-20 10:00:00+08'),
(10, 23, 'claude-sonnet-4-20250514', 65000, 25000, 1.250000, '2025-12-26 16:00:00+08'),
(2, 24,  'claude-sonnet-4-20250514', 40000, 16000, 0.750000, '2026-01-02 14:00:00+08'),

-- Idea 7: 健身饮食
(6, 25, 'claude-sonnet-4-20250514', 12000, 5000, 0.150000, '2026-03-03 11:00:00+08');

-- ============================================================
-- Contributions (for approved tasks)
-- ============================================================
INSERT INTO contributions (idea_id, task_id, user_id, cost_usd, quality_score, weighted_score, weight_percent) VALUES
-- Idea 1 (doc1 and doc2 approved)
(1, 1, 2, 0.850000, 1.20, 1.020000, 0),
(1, 2, 3, 0.620000, 1.50, 0.930000, 0),

-- Idea 2 (doc1 approved)
(2, 5, 4, 0.720000, 1.00, 0.720000, 100.0000),

-- Idea 3 (all approved, weights calculated)
(3, 9,  2, 0.900000, 1.20, 1.080000, 27.9000),
(3, 10, 6, 0.550000, 1.50, 0.825000, 21.3000),
(3, 11, 2, 1.100000, 1.20, 1.320000, 34.1000),
(3, 12, 7, 0.650000, 1.00, 0.650000, 16.7000),

-- Idea 5 (doc1 and doc2 approved)
(5, 17, 3, 0.880000, 1.20, 1.056000, 0),
(5, 18, 8, 0.590000, 1.50, 0.885000, 0),

-- Idea 6 (all approved, weights calculated)
(6, 21, 4,  0.780000, 1.00, 0.780000, 20.2000),
(6, 22, 7,  0.520000, 1.20, 0.624000, 16.1000),
(6, 23, 10, 1.250000, 1.50, 1.875000, 48.5000),
(6, 24, 2,  0.750000, 1.20, 0.900000, 15.2000);

-- ============================================================
-- PRDs for completed ideas (2 PRDs)
-- ============================================================
INSERT INTO prds (id, idea_id, content, published_at, price_credits, read_count) VALUES
(1, 3,
 E'# PetCare+ — 宠物健康管理 App 完整产品需求文档\n\n帮助宠物主人记录宠物的饮食、运动、疫苗、体检等健康数据，AI 分析异常指标并提醒就医。\n\n---\n\n# doc1 - 竞品分析报告\n\n## 直接竞品\n1. 小壹管家 — 国内宠物综合管理\n2. PetDesk — 美国市场兽医预约\n3. 11pets — 欧洲健康记录\n\n---\n\n# doc2 - 目标用户画像\n\n新手铲屎官、多宠家庭、宠物店主\n\n---\n\n# doc3 - 产品需求文档\n\n## 核心功能\n1. 宠物档案管理（多宠物）\n2. 健康数据记录\n3. 疫苗提醒\n4. AI 异常预警\n5. 兽医问诊入口\n\n## 商业模式\n- 基础版免费\n- Pro 版 ¥25/月\n\n---\n\n# doc4 - 技术可行性评估\n\n可行。React Native + Go + PostgreSQL。MVP 3 个月。',
 '2026-01-25 09:00:00+08', 6400.0000, 5),

(2, 6,
 E'# LiveClass — 在线教育直播平台完整产品需求文档\n\n面向 K12 和成人教育的互动直播教学平台。\n\n---\n\n# doc1 - 竞品分析报告\n\n腾讯课堂、ClassIn、钉钉课堂分析。差异化: AI 助教 + 低价。\n\n---\n\n# doc2 - 目标用户画像\n\nK12 数学老师、独立培训讲师。\n\n---\n\n# doc3 - 产品需求文档\n\n直播间管理、白板协作、弹幕答题、AI 课堂总结、学生签到。\n\n---\n\n# doc4 - 技术可行性评估\n\n可行。Next.js + Go + 声网 Agora SDK。',
 '2026-01-05 09:00:00+08', 6600.0000, 8);

SELECT setval('prds_id_seq', 2);

-- ============================================================
-- Credit transactions (diverse transaction types)
-- ============================================================
INSERT INTO credit_transactions (user_id, type, amount, reference_type, reference_id, description, created_at) VALUES
-- Contribution rewards (Idea 1)
(2, 'earn_contribute', 1020.0000, 'task', 1, 'doc1 竞品分析 — AI 笔记助手 (quality: 1.2)', '2026-01-21 10:00:00+08'),
(3, 'earn_contribute', 930.0000, 'task', 2, 'doc2 用户画像 — AI 笔记助手 (quality: 1.5)', '2026-01-23 09:00:00+08'),

-- Contribution rewards (Idea 2)
(4, 'earn_contribute', 720.0000, 'task', 5, 'doc1 竞品分析 — 社区团购 (quality: 1.0)', '2026-02-01 09:00:00+08'),

-- Contribution rewards (Idea 3 - all 4 tasks)
(2, 'earn_contribute', 1080.0000, 'task', 9,  'doc1 竞品分析 — 宠物健康管理 (quality: 1.2)', '2026-01-09 09:00:00+08'),
(6, 'earn_contribute', 825.0000,  'task', 10, 'doc2 用户画像 — 宠物健康管理 (quality: 1.5)', '2026-01-11 10:00:00+08'),
(2, 'earn_contribute', 1320.0000, 'task', 11, 'doc3 产品需求文档 — 宠物健康管理 (quality: 1.2)', '2026-01-17 11:00:00+08'),
(7, 'earn_contribute', 650.0000,  'task', 12, 'doc4 技术可行性 — 宠物健康管理 (quality: 1.0)', '2026-01-23 09:00:00+08'),

-- Contribution rewards (Idea 5)
(3, 'earn_contribute', 1056.0000, 'task', 17, 'doc1 竞品分析 — 商标交易 (quality: 1.2)', '2026-02-09 09:00:00+08'),
(8, 'earn_contribute', 885.0000,  'task', 18, 'doc2 用户画像 — 商标交易 (quality: 1.5)', '2026-02-10 10:00:00+08'),

-- Contribution rewards (Idea 6 - all 4 tasks)
(4,  'earn_contribute', 780.0000,  'task', 21, 'doc1 竞品分析 — 在线教育 (quality: 1.0)', '2025-12-19 09:00:00+08'),
(7,  'earn_contribute', 624.0000,  'task', 22, 'doc2 用户画像 — 在线教育 (quality: 1.2)', '2025-12-21 10:00:00+08'),
(10, 'earn_contribute', 1875.0000, 'task', 23, 'doc3 产品需求文档 — 在线教育 (quality: 1.5)', '2025-12-27 11:00:00+08'),
(2,  'earn_contribute', 900.0000,  'task', 24, 'doc4 技术可行性 — 在线教育 (quality: 1.2)', '2026-01-03 09:00:00+08'),

-- PRD purchases
(5, 'spend_read', -6400.0000, 'prd', 1, '购买 PetCare+ PRD', '2026-02-05 14:00:00+08'),
(9, 'spend_read', -6400.0000, 'prd', 1, '购买 PetCare+ PRD', '2026-02-08 10:00:00+08'),
(10, 'spend_read', -6400.0000, 'prd', 1, '购买 PetCare+ PRD', '2026-02-10 09:00:00+08'),

-- PRD purchase revenue distribution (Idea 3, buyer: user 5)
(3,  'earn_initiator_cut', 1600.0000, 'prd', 1, 'PetCare+ PRD 发起人分成 (25%)', '2026-02-05 14:00:00+08'),
(2,  'earn_read_share', 1785.6000, 'prd', 1, 'PetCare+ PRD 贡献者分成', '2026-02-05 14:00:00+08'),
(6,  'earn_read_share', 1364.8000, 'prd', 1, 'PetCare+ PRD 贡献者分成', '2026-02-05 14:00:00+08'),
(7,  'earn_read_share',  681.6000, 'prd', 1, 'PetCare+ PRD 贡献者分成', '2026-02-05 14:00:00+08'),

-- LiveClass PRD purchases
(1,  'spend_read', -6600.0000, 'prd', 2, '购买 LiveClass PRD', '2026-01-10 09:00:00+08'),
(3,  'spend_read', -6600.0000, 'prd', 2, '购买 LiveClass PRD', '2026-01-15 11:00:00+08'),

-- LiveClass PRD revenue distribution
(6,  'earn_initiator_cut', 1320.0000, 'prd', 2, 'LiveClass PRD 发起人分成 (20%)', '2026-01-10 09:00:00+08'),
(4,  'earn_read_share', 933.2400, 'prd', 2, 'LiveClass PRD 贡献者分成', '2026-01-10 09:00:00+08'),
(7,  'earn_read_share', 743.8200, 'prd', 2, 'LiveClass PRD 贡献者分成', '2026-01-10 09:00:00+08'),
(10, 'earn_read_share', 2240.1000, 'prd', 2, 'LiveClass PRD 贡献者分成', '2026-01-10 09:00:00+08'),
(2,  'earn_read_share', 702.2400, 'prd', 2, 'LiveClass PRD 贡献者分成', '2026-01-10 09:00:00+08');

-- ============================================================
-- OAuth accounts (X login)
-- ============================================================
INSERT INTO user_oauth_accounts (user_id, provider, provider_user_id, provider_username) VALUES
(1,  'x', '1001', 'alice_chen'),
(2,  'x', '1002', 'bob_builds'),
(3,  'x', '1003', 'carol_designs'),
(4,  'x', '1004', 'david_codes'),
(5,  'x', '1005', 'eve_creates'),
(6,  'x', '1006', 'frank_dev'),
(7,  'x', '1007', 'grace_pm'),
(8,  'x', '1008', 'henry_writes'),
(9,  'x', '1009', 'iris_data'),
(10, 'x', '1010', 'jack_full');

COMMIT;
