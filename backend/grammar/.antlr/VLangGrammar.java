// Generated from c:/Users/72358/Desktop/Compiler/OLC2_Proyecto1_201905884/backend/grammar/VLangGrammar.g4 by ANTLR 4.13.1
import org.antlr.v4.runtime.atn.*;
import org.antlr.v4.runtime.dfa.DFA;
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.misc.*;
import org.antlr.v4.runtime.tree.*;
import java.util.List;
import java.util.Iterator;
import java.util.ArrayList;

@SuppressWarnings({"all", "warnings", "unchecked", "unused", "cast", "CheckReturnValue"})
public class VLangGrammar extends Parser {
	static { RuntimeMetaData.checkVersion("4.13.1", RuntimeMetaData.VERSION); }

	protected static final DFA[] _decisionToDFA;
	protected static final PredictionContextCache _sharedContextCache =
		new PredictionContextCache();
	public static final int
		MUT=1, FUNC=2, STR=3, IF_KW=4, ELSE_KW=5, SWITCH_KW=6, CASE_KW=7, DEFAULT_KW=8, 
		FOR_KW=9, WHILE_KW=10, IN_KW=11, BREAK_KW=12, CONTINUE_KW=13, RETURN_KW=14, 
		DEC=15, INC=16, PLUS=17, MINUS=18, MULT=19, DIV=20, MOD=21, ASSIGN=22, 
		PLUS_ASSIGN=23, MINUS_ASSIGN=24, EQ=25, NE=26, LT=27, LE=28, GT=29, GE=30, 
		AND=31, OR=32, NOT=33, LPAREN=34, RPAREN=35, LBRACE=36, RBRACE=37, LBRACK=38, 
		RBRACK=39, SEMI=40, COLON=41, DOT=42, COMMA=43, INT_LITERAL=44, FLOAT_LITERAL=45, 
		STRING_LITERAL=46, BOOL_LITERAL=47, NIL_LITERAL=48, ID=49, WS=50, LINE_COMMENT=51, 
		BLOCK_COMMENT=52;
	public static final int
		RULE_program = 0, RULE_stmt = 1, RULE_decl_stmt = 2, RULE_var_type = 3, 
		RULE_vect_expr = 4, RULE_vect_item = 5, RULE_vect_prop = 6, RULE_vect_func = 7, 
		RULE_repeating = 8, RULE_vector_type = 9, RULE_matrix_type = 10, RULE_aux_matrix_type = 11, 
		RULE_type = 12, RULE_assign_stmt = 13, RULE_id_pattern = 14, RULE_literal = 15, 
		RULE_incredecre = 16, RULE_expression = 17, RULE_if_stmt = 18, RULE_if_chain = 19, 
		RULE_else_stmt = 20, RULE_switch_stmt = 21, RULE_switch_case = 22, RULE_default_case = 23, 
		RULE_while_stmt = 24, RULE_for_stmt = 25, RULE_range = 26, RULE_transfer_stmt = 27, 
		RULE_func_call = 28, RULE_block_ind = 29, RULE_arg_list = 30, RULE_func_arg = 31, 
		RULE_func_dcl = 32, RULE_param_list = 33, RULE_func_param = 34, RULE_strct_dcl = 35, 
		RULE_struct_prop = 36, RULE_struct_param_list = 37, RULE_struct_param = 38;
	private static String[] makeRuleNames() {
		return new String[] {
			"program", "stmt", "decl_stmt", "var_type", "vect_expr", "vect_item", 
			"vect_prop", "vect_func", "repeating", "vector_type", "matrix_type", 
			"aux_matrix_type", "type", "assign_stmt", "id_pattern", "literal", "incredecre", 
			"expression", "if_stmt", "if_chain", "else_stmt", "switch_stmt", "switch_case", 
			"default_case", "while_stmt", "for_stmt", "range", "transfer_stmt", "func_call", 
			"block_ind", "arg_list", "func_arg", "func_dcl", "param_list", "func_param", 
			"strct_dcl", "struct_prop", "struct_param_list", "struct_param"
		};
	}
	public static final String[] ruleNames = makeRuleNames();

	private static String[] makeLiteralNames() {
		return new String[] {
			null, "'mut'", "'fn'", "'struct'", "'if'", "'else'", "'switch'", "'case'", 
			"'default'", "'for'", "'while'", "'in'", "'break'", "'continue'", "'return'", 
			"'--'", "'++'", "'+'", "'-'", "'*'", "'/'", "'%'", "'='", "'+='", "'-='", 
			"'=='", "'!='", "'<'", "'<='", "'>'", "'>='", "'&&'", "'||'", "'!'", 
			"'('", "')'", "'{'", "'}'", "'['", "']'", "';'", "':'", "'.'", "','", 
			null, null, null, null, "'nil'"
		};
	}
	private static final String[] _LITERAL_NAMES = makeLiteralNames();
	private static String[] makeSymbolicNames() {
		return new String[] {
			null, "MUT", "FUNC", "STR", "IF_KW", "ELSE_KW", "SWITCH_KW", "CASE_KW", 
			"DEFAULT_KW", "FOR_KW", "WHILE_KW", "IN_KW", "BREAK_KW", "CONTINUE_KW", 
			"RETURN_KW", "DEC", "INC", "PLUS", "MINUS", "MULT", "DIV", "MOD", "ASSIGN", 
			"PLUS_ASSIGN", "MINUS_ASSIGN", "EQ", "NE", "LT", "LE", "GT", "GE", "AND", 
			"OR", "NOT", "LPAREN", "RPAREN", "LBRACE", "RBRACE", "LBRACK", "RBRACK", 
			"SEMI", "COLON", "DOT", "COMMA", "INT_LITERAL", "FLOAT_LITERAL", "STRING_LITERAL", 
			"BOOL_LITERAL", "NIL_LITERAL", "ID", "WS", "LINE_COMMENT", "BLOCK_COMMENT"
		};
	}
	private static final String[] _SYMBOLIC_NAMES = makeSymbolicNames();
	public static final Vocabulary VOCABULARY = new VocabularyImpl(_LITERAL_NAMES, _SYMBOLIC_NAMES);

	/**
	 * @deprecated Use {@link #VOCABULARY} instead.
	 */
	@Deprecated
	public static final String[] tokenNames;
	static {
		tokenNames = new String[_SYMBOLIC_NAMES.length];
		for (int i = 0; i < tokenNames.length; i++) {
			tokenNames[i] = VOCABULARY.getLiteralName(i);
			if (tokenNames[i] == null) {
				tokenNames[i] = VOCABULARY.getSymbolicName(i);
			}

			if (tokenNames[i] == null) {
				tokenNames[i] = "<INVALID>";
			}
		}
	}

	@Override
	@Deprecated
	public String[] getTokenNames() {
		return tokenNames;
	}

	@Override

	public Vocabulary getVocabulary() {
		return VOCABULARY;
	}

	@Override
	public String getGrammarFileName() { return "VLangGrammar.g4"; }

	@Override
	public String[] getRuleNames() { return ruleNames; }

	@Override
	public String getSerializedATN() { return _serializedATN; }

	@Override
	public ATN getATN() { return _ATN; }

	public VLangGrammar(TokenStream input) {
		super(input);
		_interp = new ParserATNSimulator(this,_ATN,_decisionToDFA,_sharedContextCache);
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ProgramContext extends ParserRuleContext {
		public List<StmtContext> stmt() {
			return getRuleContexts(StmtContext.class);
		}
		public StmtContext stmt(int i) {
			return getRuleContext(StmtContext.class,i);
		}
		public TerminalNode EOF() { return getToken(VLangGrammar.EOF, 0); }
		public ProgramContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_program; }
	}

	public final ProgramContext program() throws RecognitionException {
		ProgramContext _localctx = new ProgramContext(_ctx, getState());
		enterRule(_localctx, 0, RULE_program);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(81);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 563018672928350L) != 0)) {
				{
				{
				setState(78);
				stmt();
				}
				}
				setState(83);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(85);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,1,_ctx) ) {
			case 1:
				{
				setState(84);
				match(EOF);
				}
				break;
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class StmtContext extends ParserRuleContext {
		public Decl_stmtContext decl_stmt() {
			return getRuleContext(Decl_stmtContext.class,0);
		}
		public Assign_stmtContext assign_stmt() {
			return getRuleContext(Assign_stmtContext.class,0);
		}
		public Block_indContext block_ind() {
			return getRuleContext(Block_indContext.class,0);
		}
		public Transfer_stmtContext transfer_stmt() {
			return getRuleContext(Transfer_stmtContext.class,0);
		}
		public If_stmtContext if_stmt() {
			return getRuleContext(If_stmtContext.class,0);
		}
		public Switch_stmtContext switch_stmt() {
			return getRuleContext(Switch_stmtContext.class,0);
		}
		public While_stmtContext while_stmt() {
			return getRuleContext(While_stmtContext.class,0);
		}
		public For_stmtContext for_stmt() {
			return getRuleContext(For_stmtContext.class,0);
		}
		public Func_callContext func_call() {
			return getRuleContext(Func_callContext.class,0);
		}
		public Vect_funcContext vect_func() {
			return getRuleContext(Vect_funcContext.class,0);
		}
		public Func_dclContext func_dcl() {
			return getRuleContext(Func_dclContext.class,0);
		}
		public Strct_dclContext strct_dcl() {
			return getRuleContext(Strct_dclContext.class,0);
		}
		public StmtContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_stmt; }
	}

	public final StmtContext stmt() throws RecognitionException {
		StmtContext _localctx = new StmtContext(_ctx, getState());
		enterRule(_localctx, 2, RULE_stmt);
		try {
			setState(99);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,2,_ctx) ) {
			case 1:
				enterOuterAlt(_localctx, 1);
				{
				setState(87);
				decl_stmt();
				}
				break;
			case 2:
				enterOuterAlt(_localctx, 2);
				{
				setState(88);
				assign_stmt();
				}
				break;
			case 3:
				enterOuterAlt(_localctx, 3);
				{
				setState(89);
				block_ind();
				}
				break;
			case 4:
				enterOuterAlt(_localctx, 4);
				{
				setState(90);
				transfer_stmt();
				}
				break;
			case 5:
				enterOuterAlt(_localctx, 5);
				{
				setState(91);
				if_stmt();
				}
				break;
			case 6:
				enterOuterAlt(_localctx, 6);
				{
				setState(92);
				switch_stmt();
				}
				break;
			case 7:
				enterOuterAlt(_localctx, 7);
				{
				setState(93);
				while_stmt();
				}
				break;
			case 8:
				enterOuterAlt(_localctx, 8);
				{
				setState(94);
				for_stmt();
				}
				break;
			case 9:
				enterOuterAlt(_localctx, 9);
				{
				setState(95);
				func_call();
				}
				break;
			case 10:
				enterOuterAlt(_localctx, 10);
				{
				setState(96);
				vect_func();
				}
				break;
			case 11:
				enterOuterAlt(_localctx, 11);
				{
				setState(97);
				func_dcl();
				}
				break;
			case 12:
				enterOuterAlt(_localctx, 12);
				{
				setState(98);
				strct_dcl();
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Decl_stmtContext extends ParserRuleContext {
		public Decl_stmtContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_decl_stmt; }
	 
		public Decl_stmtContext() { }
		public void copyFrom(Decl_stmtContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class VarAssDeclContext extends Decl_stmtContext {
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public TypeContext type() {
			return getRuleContext(TypeContext.class,0);
		}
		public TerminalNode ASSIGN() { return getToken(VLangGrammar.ASSIGN, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public VarAssDeclContext(Decl_stmtContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class VarMatrixDeclContext extends Decl_stmtContext {
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public TerminalNode ASSIGN() { return getToken(VLangGrammar.ASSIGN, 0); }
		public Matrix_typeContext matrix_type() {
			return getRuleContext(Matrix_typeContext.class,0);
		}
		public Vect_exprContext vect_expr() {
			return getRuleContext(Vect_exprContext.class,0);
		}
		public VarMatrixDeclContext(Decl_stmtContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class ValDeclVecContext extends Decl_stmtContext {
		public Var_typeContext var_type() {
			return getRuleContext(Var_typeContext.class,0);
		}
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public TypeContext type() {
			return getRuleContext(TypeContext.class,0);
		}
		public ValDeclVecContext(Decl_stmtContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class ValueDeclContext extends Decl_stmtContext {
		public Var_typeContext var_type() {
			return getRuleContext(Var_typeContext.class,0);
		}
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public TerminalNode ASSIGN() { return getToken(VLangGrammar.ASSIGN, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public ValueDeclContext(Decl_stmtContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class MutVarDeclContext extends Decl_stmtContext {
		public Var_typeContext var_type() {
			return getRuleContext(Var_typeContext.class,0);
		}
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public TypeContext type() {
			return getRuleContext(TypeContext.class,0);
		}
		public TerminalNode ASSIGN() { return getToken(VLangGrammar.ASSIGN, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public MutVarDeclContext(Decl_stmtContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class VarVectDeclContext extends Decl_stmtContext {
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public TerminalNode ASSIGN() { return getToken(VLangGrammar.ASSIGN, 0); }
		public Vector_typeContext vector_type() {
			return getRuleContext(Vector_typeContext.class,0);
		}
		public Vect_exprContext vect_expr() {
			return getRuleContext(Vect_exprContext.class,0);
		}
		public VarVectDeclContext(Decl_stmtContext ctx) { copyFrom(ctx); }
	}

	public final Decl_stmtContext decl_stmt() throws RecognitionException {
		Decl_stmtContext _localctx = new Decl_stmtContext(_ctx, getState());
		enterRule(_localctx, 4, RULE_decl_stmt);
		try {
			setState(131);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,3,_ctx) ) {
			case 1:
				_localctx = new MutVarDeclContext(_localctx);
				enterOuterAlt(_localctx, 1);
				{
				setState(101);
				var_type();
				setState(102);
				match(ID);
				setState(103);
				type();
				setState(104);
				match(ASSIGN);
				setState(105);
				expression(0);
				}
				break;
			case 2:
				_localctx = new ValueDeclContext(_localctx);
				enterOuterAlt(_localctx, 2);
				{
				setState(107);
				var_type();
				setState(108);
				match(ID);
				setState(109);
				match(ASSIGN);
				setState(110);
				expression(0);
				}
				break;
			case 3:
				_localctx = new ValDeclVecContext(_localctx);
				enterOuterAlt(_localctx, 3);
				{
				setState(112);
				var_type();
				setState(113);
				match(ID);
				setState(114);
				type();
				}
				break;
			case 4:
				_localctx = new VarAssDeclContext(_localctx);
				enterOuterAlt(_localctx, 4);
				{
				setState(116);
				match(ID);
				setState(117);
				type();
				setState(118);
				match(ASSIGN);
				setState(119);
				expression(0);
				}
				break;
			case 5:
				_localctx = new VarVectDeclContext(_localctx);
				enterOuterAlt(_localctx, 5);
				{
				setState(121);
				match(ID);
				setState(122);
				match(ASSIGN);
				setState(123);
				vector_type();
				setState(124);
				vect_expr();
				}
				break;
			case 6:
				_localctx = new VarMatrixDeclContext(_localctx);
				enterOuterAlt(_localctx, 6);
				{
				setState(126);
				match(ID);
				setState(127);
				match(ASSIGN);
				setState(128);
				matrix_type();
				setState(129);
				vect_expr();
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Var_typeContext extends ParserRuleContext {
		public TerminalNode MUT() { return getToken(VLangGrammar.MUT, 0); }
		public Var_typeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_var_type; }
	}

	public final Var_typeContext var_type() throws RecognitionException {
		Var_typeContext _localctx = new Var_typeContext(_ctx, getState());
		enterRule(_localctx, 6, RULE_var_type);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(133);
			match(MUT);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Vect_exprContext extends ParserRuleContext {
		public Vect_exprContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_vect_expr; }
	 
		public Vect_exprContext() { }
		public void copyFrom(Vect_exprContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class VectorItemLisContext extends Vect_exprContext {
		public TerminalNode LBRACE() { return getToken(VLangGrammar.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(VLangGrammar.RBRACE, 0); }
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(VLangGrammar.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(VLangGrammar.COMMA, i);
		}
		public VectorItemLisContext(Vect_exprContext ctx) { copyFrom(ctx); }
	}

	public final Vect_exprContext vect_expr() throws RecognitionException {
		Vect_exprContext _localctx = new Vect_exprContext(_ctx, getState());
		enterRule(_localctx, 8, RULE_vect_expr);
		int _la;
		try {
			_localctx = new VectorItemLisContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(135);
			match(LBRACE);
			setState(144);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if ((((_la) & ~0x3f) == 0 && ((1L << _la) & 1108677088247808L) != 0)) {
				{
				setState(136);
				expression(0);
				setState(141);
				_errHandler.sync(this);
				_la = _input.LA(1);
				while (_la==COMMA) {
					{
					{
					setState(137);
					match(COMMA);
					setState(138);
					expression(0);
					}
					}
					setState(143);
					_errHandler.sync(this);
					_la = _input.LA(1);
				}
				}
			}

			setState(146);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Vect_itemContext extends ParserRuleContext {
		public Vect_itemContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_vect_item; }
	 
		public Vect_itemContext() { }
		public void copyFrom(Vect_itemContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class VectorItemContext extends Vect_itemContext {
		public Id_patternContext id_pattern() {
			return getRuleContext(Id_patternContext.class,0);
		}
		public List<TerminalNode> LBRACK() { return getTokens(VLangGrammar.LBRACK); }
		public TerminalNode LBRACK(int i) {
			return getToken(VLangGrammar.LBRACK, i);
		}
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public List<TerminalNode> RBRACK() { return getTokens(VLangGrammar.RBRACK); }
		public TerminalNode RBRACK(int i) {
			return getToken(VLangGrammar.RBRACK, i);
		}
		public VectorItemContext(Vect_itemContext ctx) { copyFrom(ctx); }
	}

	public final Vect_itemContext vect_item() throws RecognitionException {
		Vect_itemContext _localctx = new Vect_itemContext(_ctx, getState());
		enterRule(_localctx, 10, RULE_vect_item);
		try {
			int _alt;
			_localctx = new VectorItemContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(148);
			id_pattern();
			setState(153); 
			_errHandler.sync(this);
			_alt = 1;
			do {
				switch (_alt) {
				case 1:
					{
					{
					setState(149);
					match(LBRACK);
					setState(150);
					expression(0);
					setState(151);
					match(RBRACK);
					}
					}
					break;
				default:
					throw new NoViableAltException(this);
				}
				setState(155); 
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,6,_ctx);
			} while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER );
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Vect_propContext extends ParserRuleContext {
		public Vect_propContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_vect_prop; }
	 
		public Vect_propContext() { }
		public void copyFrom(Vect_propContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class VectorPropertyContext extends Vect_propContext {
		public Vect_itemContext vect_item() {
			return getRuleContext(Vect_itemContext.class,0);
		}
		public TerminalNode DOT() { return getToken(VLangGrammar.DOT, 0); }
		public Id_patternContext id_pattern() {
			return getRuleContext(Id_patternContext.class,0);
		}
		public VectorPropertyContext(Vect_propContext ctx) { copyFrom(ctx); }
	}

	public final Vect_propContext vect_prop() throws RecognitionException {
		Vect_propContext _localctx = new Vect_propContext(_ctx, getState());
		enterRule(_localctx, 12, RULE_vect_prop);
		try {
			_localctx = new VectorPropertyContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(157);
			vect_item();
			setState(158);
			match(DOT);
			setState(159);
			id_pattern();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Vect_funcContext extends ParserRuleContext {
		public Vect_funcContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_vect_func; }
	 
		public Vect_funcContext() { }
		public void copyFrom(Vect_funcContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class VectorFuncCallContext extends Vect_funcContext {
		public Vect_itemContext vect_item() {
			return getRuleContext(Vect_itemContext.class,0);
		}
		public TerminalNode DOT() { return getToken(VLangGrammar.DOT, 0); }
		public Func_callContext func_call() {
			return getRuleContext(Func_callContext.class,0);
		}
		public VectorFuncCallContext(Vect_funcContext ctx) { copyFrom(ctx); }
	}

	public final Vect_funcContext vect_func() throws RecognitionException {
		Vect_funcContext _localctx = new Vect_funcContext(_ctx, getState());
		enterRule(_localctx, 14, RULE_vect_func);
		try {
			_localctx = new VectorFuncCallContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(161);
			vect_item();
			setState(162);
			match(DOT);
			setState(163);
			func_call();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class RepeatingContext extends ParserRuleContext {
		public RepeatingContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_repeating; }
	 
		public RepeatingContext() { }
		public void copyFrom(RepeatingContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class RepeatingDeclContext extends RepeatingContext {
		public TerminalNode LPAREN() { return getToken(VLangGrammar.LPAREN, 0); }
		public List<TerminalNode> ID() { return getTokens(VLangGrammar.ID); }
		public TerminalNode ID(int i) {
			return getToken(VLangGrammar.ID, i);
		}
		public List<TerminalNode> COLON() { return getTokens(VLangGrammar.COLON); }
		public TerminalNode COLON(int i) {
			return getToken(VLangGrammar.COLON, i);
		}
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode COMMA() { return getToken(VLangGrammar.COMMA, 0); }
		public TerminalNode RPAREN() { return getToken(VLangGrammar.RPAREN, 0); }
		public Vector_typeContext vector_type() {
			return getRuleContext(Vector_typeContext.class,0);
		}
		public Matrix_typeContext matrix_type() {
			return getRuleContext(Matrix_typeContext.class,0);
		}
		public RepeatingDeclContext(RepeatingContext ctx) { copyFrom(ctx); }
	}

	public final RepeatingContext repeating() throws RecognitionException {
		RepeatingContext _localctx = new RepeatingContext(_ctx, getState());
		enterRule(_localctx, 16, RULE_repeating);
		try {
			_localctx = new RepeatingDeclContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(167);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,7,_ctx) ) {
			case 1:
				{
				setState(165);
				vector_type();
				}
				break;
			case 2:
				{
				setState(166);
				matrix_type();
				}
				break;
			}
			setState(169);
			match(LPAREN);
			setState(170);
			match(ID);
			setState(171);
			match(COLON);
			setState(172);
			expression(0);
			setState(173);
			match(COMMA);
			setState(174);
			match(ID);
			setState(175);
			match(COLON);
			setState(176);
			expression(0);
			setState(177);
			match(RPAREN);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Vector_typeContext extends ParserRuleContext {
		public TerminalNode LBRACK() { return getToken(VLangGrammar.LBRACK, 0); }
		public TerminalNode RBRACK() { return getToken(VLangGrammar.RBRACK, 0); }
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public Vector_typeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_vector_type; }
	}

	public final Vector_typeContext vector_type() throws RecognitionException {
		Vector_typeContext _localctx = new Vector_typeContext(_ctx, getState());
		enterRule(_localctx, 18, RULE_vector_type);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(179);
			match(LBRACK);
			setState(180);
			match(RBRACK);
			setState(181);
			match(ID);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Matrix_typeContext extends ParserRuleContext {
		public Aux_matrix_typeContext aux_matrix_type() {
			return getRuleContext(Aux_matrix_typeContext.class,0);
		}
		public List<TerminalNode> LBRACK() { return getTokens(VLangGrammar.LBRACK); }
		public TerminalNode LBRACK(int i) {
			return getToken(VLangGrammar.LBRACK, i);
		}
		public List<TerminalNode> RBRACK() { return getTokens(VLangGrammar.RBRACK); }
		public TerminalNode RBRACK(int i) {
			return getToken(VLangGrammar.RBRACK, i);
		}
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public Matrix_typeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_matrix_type; }
	}

	public final Matrix_typeContext matrix_type() throws RecognitionException {
		Matrix_typeContext _localctx = new Matrix_typeContext(_ctx, getState());
		enterRule(_localctx, 20, RULE_matrix_type);
		try {
			setState(189);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,8,_ctx) ) {
			case 1:
				enterOuterAlt(_localctx, 1);
				{
				setState(183);
				aux_matrix_type();
				}
				break;
			case 2:
				enterOuterAlt(_localctx, 2);
				{
				setState(184);
				match(LBRACK);
				setState(185);
				match(RBRACK);
				setState(186);
				match(LBRACK);
				setState(187);
				match(RBRACK);
				setState(188);
				match(ID);
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Aux_matrix_typeContext extends ParserRuleContext {
		public TerminalNode LBRACK() { return getToken(VLangGrammar.LBRACK, 0); }
		public TerminalNode RBRACK() { return getToken(VLangGrammar.RBRACK, 0); }
		public Matrix_typeContext matrix_type() {
			return getRuleContext(Matrix_typeContext.class,0);
		}
		public Aux_matrix_typeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_aux_matrix_type; }
	}

	public final Aux_matrix_typeContext aux_matrix_type() throws RecognitionException {
		Aux_matrix_typeContext _localctx = new Aux_matrix_typeContext(_ctx, getState());
		enterRule(_localctx, 22, RULE_aux_matrix_type);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(191);
			match(LBRACK);
			setState(192);
			match(RBRACK);
			setState(193);
			matrix_type();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class TypeContext extends ParserRuleContext {
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public Vector_typeContext vector_type() {
			return getRuleContext(Vector_typeContext.class,0);
		}
		public Matrix_typeContext matrix_type() {
			return getRuleContext(Matrix_typeContext.class,0);
		}
		public TypeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_type; }
	}

	public final TypeContext type() throws RecognitionException {
		TypeContext _localctx = new TypeContext(_ctx, getState());
		enterRule(_localctx, 24, RULE_type);
		try {
			setState(198);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,9,_ctx) ) {
			case 1:
				enterOuterAlt(_localctx, 1);
				{
				setState(195);
				match(ID);
				}
				break;
			case 2:
				enterOuterAlt(_localctx, 2);
				{
				setState(196);
				vector_type();
				}
				break;
			case 3:
				enterOuterAlt(_localctx, 3);
				{
				setState(197);
				matrix_type();
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Assign_stmtContext extends ParserRuleContext {
		public Assign_stmtContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_assign_stmt; }
	 
		public Assign_stmtContext() { }
		public void copyFrom(Assign_stmtContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class ArgAddAssigDeclContext extends Assign_stmtContext {
		public Token op;
		public Id_patternContext id_pattern() {
			return getRuleContext(Id_patternContext.class,0);
		}
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode PLUS_ASSIGN() { return getToken(VLangGrammar.PLUS_ASSIGN, 0); }
		public TerminalNode MINUS_ASSIGN() { return getToken(VLangGrammar.MINUS_ASSIGN, 0); }
		public ArgAddAssigDeclContext(Assign_stmtContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class VectorAssignContext extends Assign_stmtContext {
		public Token op;
		public Vect_itemContext vect_item() {
			return getRuleContext(Vect_itemContext.class,0);
		}
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode PLUS_ASSIGN() { return getToken(VLangGrammar.PLUS_ASSIGN, 0); }
		public TerminalNode MINUS_ASSIGN() { return getToken(VLangGrammar.MINUS_ASSIGN, 0); }
		public TerminalNode ASSIGN() { return getToken(VLangGrammar.ASSIGN, 0); }
		public VectorAssignContext(Assign_stmtContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class AssignmentDeclContext extends Assign_stmtContext {
		public Id_patternContext id_pattern() {
			return getRuleContext(Id_patternContext.class,0);
		}
		public TerminalNode ASSIGN() { return getToken(VLangGrammar.ASSIGN, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public AssignmentDeclContext(Assign_stmtContext ctx) { copyFrom(ctx); }
	}

	public final Assign_stmtContext assign_stmt() throws RecognitionException {
		Assign_stmtContext _localctx = new Assign_stmtContext(_ctx, getState());
		enterRule(_localctx, 26, RULE_assign_stmt);
		int _la;
		try {
			setState(212);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,10,_ctx) ) {
			case 1:
				_localctx = new AssignmentDeclContext(_localctx);
				enterOuterAlt(_localctx, 1);
				{
				setState(200);
				id_pattern();
				setState(201);
				match(ASSIGN);
				setState(202);
				expression(0);
				}
				break;
			case 2:
				_localctx = new ArgAddAssigDeclContext(_localctx);
				enterOuterAlt(_localctx, 2);
				{
				setState(204);
				id_pattern();
				setState(205);
				((ArgAddAssigDeclContext)_localctx).op = _input.LT(1);
				_la = _input.LA(1);
				if ( !(_la==PLUS_ASSIGN || _la==MINUS_ASSIGN) ) {
					((ArgAddAssigDeclContext)_localctx).op = (Token)_errHandler.recoverInline(this);
				}
				else {
					if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
					_errHandler.reportMatch(this);
					consume();
				}
				setState(206);
				expression(0);
				}
				break;
			case 3:
				_localctx = new VectorAssignContext(_localctx);
				enterOuterAlt(_localctx, 3);
				{
				setState(208);
				vect_item();
				setState(209);
				((VectorAssignContext)_localctx).op = _input.LT(1);
				_la = _input.LA(1);
				if ( !((((_la) & ~0x3f) == 0 && ((1L << _la) & 29360128L) != 0)) ) {
					((VectorAssignContext)_localctx).op = (Token)_errHandler.recoverInline(this);
				}
				else {
					if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
					_errHandler.reportMatch(this);
					consume();
				}
				setState(210);
				expression(0);
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Id_patternContext extends ParserRuleContext {
		public Id_patternContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_id_pattern; }
	 
		public Id_patternContext() { }
		public void copyFrom(Id_patternContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class IdPatternContext extends Id_patternContext {
		public List<TerminalNode> ID() { return getTokens(VLangGrammar.ID); }
		public TerminalNode ID(int i) {
			return getToken(VLangGrammar.ID, i);
		}
		public List<TerminalNode> DOT() { return getTokens(VLangGrammar.DOT); }
		public TerminalNode DOT(int i) {
			return getToken(VLangGrammar.DOT, i);
		}
		public IdPatternContext(Id_patternContext ctx) { copyFrom(ctx); }
	}

	public final Id_patternContext id_pattern() throws RecognitionException {
		Id_patternContext _localctx = new Id_patternContext(_ctx, getState());
		enterRule(_localctx, 28, RULE_id_pattern);
		try {
			int _alt;
			_localctx = new IdPatternContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(214);
			match(ID);
			setState(219);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,11,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(215);
					match(DOT);
					setState(216);
					match(ID);
					}
					} 
				}
				setState(221);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,11,_ctx);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class LiteralContext extends ParserRuleContext {
		public LiteralContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_literal; }
	 
		public LiteralContext() { }
		public void copyFrom(LiteralContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class StringLiteralContext extends LiteralContext {
		public TerminalNode STRING_LITERAL() { return getToken(VLangGrammar.STRING_LITERAL, 0); }
		public StringLiteralContext(LiteralContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class BoolLiteralContext extends LiteralContext {
		public TerminalNode BOOL_LITERAL() { return getToken(VLangGrammar.BOOL_LITERAL, 0); }
		public BoolLiteralContext(LiteralContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class FloatLiteralContext extends LiteralContext {
		public TerminalNode FLOAT_LITERAL() { return getToken(VLangGrammar.FLOAT_LITERAL, 0); }
		public FloatLiteralContext(LiteralContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class NilLiteralContext extends LiteralContext {
		public TerminalNode NIL_LITERAL() { return getToken(VLangGrammar.NIL_LITERAL, 0); }
		public NilLiteralContext(LiteralContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class IntLiteralContext extends LiteralContext {
		public TerminalNode INT_LITERAL() { return getToken(VLangGrammar.INT_LITERAL, 0); }
		public IntLiteralContext(LiteralContext ctx) { copyFrom(ctx); }
	}

	public final LiteralContext literal() throws RecognitionException {
		LiteralContext _localctx = new LiteralContext(_ctx, getState());
		enterRule(_localctx, 30, RULE_literal);
		try {
			setState(227);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case INT_LITERAL:
				_localctx = new IntLiteralContext(_localctx);
				enterOuterAlt(_localctx, 1);
				{
				setState(222);
				match(INT_LITERAL);
				}
				break;
			case FLOAT_LITERAL:
				_localctx = new FloatLiteralContext(_localctx);
				enterOuterAlt(_localctx, 2);
				{
				setState(223);
				match(FLOAT_LITERAL);
				}
				break;
			case STRING_LITERAL:
				_localctx = new StringLiteralContext(_localctx);
				enterOuterAlt(_localctx, 3);
				{
				setState(224);
				match(STRING_LITERAL);
				}
				break;
			case BOOL_LITERAL:
				_localctx = new BoolLiteralContext(_localctx);
				enterOuterAlt(_localctx, 4);
				{
				setState(225);
				match(BOOL_LITERAL);
				}
				break;
			case NIL_LITERAL:
				_localctx = new NilLiteralContext(_localctx);
				enterOuterAlt(_localctx, 5);
				{
				setState(226);
				match(NIL_LITERAL);
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class IncredecreContext extends ParserRuleContext {
		public IncredecreContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_incredecre; }
	 
		public IncredecreContext() { }
		public void copyFrom(IncredecreContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class IncrementoContext extends IncredecreContext {
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public TerminalNode INC() { return getToken(VLangGrammar.INC, 0); }
		public IncrementoContext(IncredecreContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class DecrementoContext extends IncredecreContext {
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public TerminalNode DEC() { return getToken(VLangGrammar.DEC, 0); }
		public DecrementoContext(IncredecreContext ctx) { copyFrom(ctx); }
	}

	public final IncredecreContext incredecre() throws RecognitionException {
		IncredecreContext _localctx = new IncredecreContext(_ctx, getState());
		enterRule(_localctx, 32, RULE_incredecre);
		try {
			setState(233);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,13,_ctx) ) {
			case 1:
				_localctx = new IncrementoContext(_localctx);
				enterOuterAlt(_localctx, 1);
				{
				setState(229);
				match(ID);
				setState(230);
				match(INC);
				}
				break;
			case 2:
				_localctx = new DecrementoContext(_localctx);
				enterOuterAlt(_localctx, 2);
				{
				setState(231);
				match(ID);
				setState(232);
				match(DEC);
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ExpressionContext extends ParserRuleContext {
		public ExpressionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_expression; }
	 
		public ExpressionContext() { }
		public void copyFrom(ExpressionContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class RepeatingExprContext extends ExpressionContext {
		public RepeatingContext repeating() {
			return getRuleContext(RepeatingContext.class,0);
		}
		public RepeatingExprContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class IncredecrContext extends ExpressionContext {
		public IncredecreContext incredecre() {
			return getRuleContext(IncredecreContext.class,0);
		}
		public IncredecrContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class StructAccessExprContext extends ExpressionContext {
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode DOT() { return getToken(VLangGrammar.DOT, 0); }
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public StructAccessExprContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class BinaryExprContext extends ExpressionContext {
		public ExpressionContext left;
		public Token op;
		public ExpressionContext right;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode MULT() { return getToken(VLangGrammar.MULT, 0); }
		public TerminalNode DIV() { return getToken(VLangGrammar.DIV, 0); }
		public TerminalNode MOD() { return getToken(VLangGrammar.MOD, 0); }
		public TerminalNode PLUS() { return getToken(VLangGrammar.PLUS, 0); }
		public TerminalNode MINUS() { return getToken(VLangGrammar.MINUS, 0); }
		public TerminalNode LE() { return getToken(VLangGrammar.LE, 0); }
		public TerminalNode LT() { return getToken(VLangGrammar.LT, 0); }
		public TerminalNode GE() { return getToken(VLangGrammar.GE, 0); }
		public TerminalNode GT() { return getToken(VLangGrammar.GT, 0); }
		public TerminalNode EQ() { return getToken(VLangGrammar.EQ, 0); }
		public TerminalNode NE() { return getToken(VLangGrammar.NE, 0); }
		public TerminalNode AND() { return getToken(VLangGrammar.AND, 0); }
		public TerminalNode OR() { return getToken(VLangGrammar.OR, 0); }
		public BinaryExprContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class StructInstantiationExprContext extends ExpressionContext {
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public TerminalNode LBRACE() { return getToken(VLangGrammar.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(VLangGrammar.RBRACE, 0); }
		public Struct_param_listContext struct_param_list() {
			return getRuleContext(Struct_param_listContext.class,0);
		}
		public StructInstantiationExprContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class UnaryExprContext extends ExpressionContext {
		public Token op;
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode NOT() { return getToken(VLangGrammar.NOT, 0); }
		public TerminalNode MINUS() { return getToken(VLangGrammar.MINUS, 0); }
		public UnaryExprContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class IdPatternExprContext extends ExpressionContext {
		public Id_patternContext id_pattern() {
			return getRuleContext(Id_patternContext.class,0);
		}
		public IdPatternExprContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class VectorPropertyExprContext extends ExpressionContext {
		public Vect_propContext vect_prop() {
			return getRuleContext(Vect_propContext.class,0);
		}
		public VectorPropertyExprContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class VectorItemExprContext extends ExpressionContext {
		public Vect_itemContext vect_item() {
			return getRuleContext(Vect_itemContext.class,0);
		}
		public VectorItemExprContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class ParensExprContext extends ExpressionContext {
		public TerminalNode LPAREN() { return getToken(VLangGrammar.LPAREN, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode RPAREN() { return getToken(VLangGrammar.RPAREN, 0); }
		public ParensExprContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class LiteralExprContext extends ExpressionContext {
		public LiteralContext literal() {
			return getRuleContext(LiteralContext.class,0);
		}
		public LiteralExprContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class VectorFuncCallExprContext extends ExpressionContext {
		public Vect_funcContext vect_func() {
			return getRuleContext(Vect_funcContext.class,0);
		}
		public VectorFuncCallExprContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class VectorExprContext extends ExpressionContext {
		public Vect_exprContext vect_expr() {
			return getRuleContext(Vect_exprContext.class,0);
		}
		public VectorExprContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class FuncCallExprContext extends ExpressionContext {
		public Func_callContext func_call() {
			return getRuleContext(Func_callContext.class,0);
		}
		public FuncCallExprContext(ExpressionContext ctx) { copyFrom(ctx); }
	}

	public final ExpressionContext expression() throws RecognitionException {
		return expression(0);
	}

	private ExpressionContext expression(int _p) throws RecognitionException {
		ParserRuleContext _parentctx = _ctx;
		int _parentState = getState();
		ExpressionContext _localctx = new ExpressionContext(_ctx, _parentState);
		ExpressionContext _prevctx = _localctx;
		int _startState = 34;
		enterRecursionRule(_localctx, 34, RULE_expression, _p);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(257);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,15,_ctx) ) {
			case 1:
				{
				_localctx = new ParensExprContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;

				setState(236);
				match(LPAREN);
				setState(237);
				expression(0);
				setState(238);
				match(RPAREN);
				}
				break;
			case 2:
				{
				_localctx = new FuncCallExprContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(240);
				func_call();
				}
				break;
			case 3:
				{
				_localctx = new IdPatternExprContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(241);
				id_pattern();
				}
				break;
			case 4:
				{
				_localctx = new VectorItemExprContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(242);
				vect_item();
				}
				break;
			case 5:
				{
				_localctx = new VectorPropertyExprContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(243);
				vect_prop();
				}
				break;
			case 6:
				{
				_localctx = new VectorFuncCallExprContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(244);
				vect_func();
				}
				break;
			case 7:
				{
				_localctx = new LiteralExprContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(245);
				literal();
				}
				break;
			case 8:
				{
				_localctx = new VectorExprContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(246);
				vect_expr();
				}
				break;
			case 9:
				{
				_localctx = new RepeatingExprContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(247);
				repeating();
				}
				break;
			case 10:
				{
				_localctx = new IncredecrContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(248);
				incredecre();
				}
				break;
			case 11:
				{
				_localctx = new UnaryExprContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(249);
				((UnaryExprContext)_localctx).op = _input.LT(1);
				_la = _input.LA(1);
				if ( !(_la==MINUS || _la==NOT) ) {
					((UnaryExprContext)_localctx).op = (Token)_errHandler.recoverInline(this);
				}
				else {
					if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
					_errHandler.reportMatch(this);
					consume();
				}
				setState(250);
				expression(9);
				}
				break;
			case 12:
				{
				_localctx = new StructInstantiationExprContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(251);
				match(ID);
				setState(252);
				match(LBRACE);
				setState(254);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==ID) {
					{
					setState(253);
					struct_param_list();
					}
				}

				setState(256);
				match(RBRACE);
				}
				break;
			}
			_ctx.stop = _input.LT(-1);
			setState(282);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,17,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					if ( _parseListeners!=null ) triggerExitRuleEvent();
					_prevctx = _localctx;
					{
					setState(280);
					_errHandler.sync(this);
					switch ( getInterpreter().adaptivePredict(_input,16,_ctx) ) {
					case 1:
						{
						_localctx = new BinaryExprContext(new ExpressionContext(_parentctx, _parentState));
						((BinaryExprContext)_localctx).left = _prevctx;
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(259);
						if (!(precpred(_ctx, 8))) throw new FailedPredicateException(this, "precpred(_ctx, 8)");
						setState(260);
						((BinaryExprContext)_localctx).op = _input.LT(1);
						_la = _input.LA(1);
						if ( !((((_la) & ~0x3f) == 0 && ((1L << _la) & 3670016L) != 0)) ) {
							((BinaryExprContext)_localctx).op = (Token)_errHandler.recoverInline(this);
						}
						else {
							if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
							_errHandler.reportMatch(this);
							consume();
						}
						setState(261);
						((BinaryExprContext)_localctx).right = expression(9);
						}
						break;
					case 2:
						{
						_localctx = new BinaryExprContext(new ExpressionContext(_parentctx, _parentState));
						((BinaryExprContext)_localctx).left = _prevctx;
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(262);
						if (!(precpred(_ctx, 7))) throw new FailedPredicateException(this, "precpred(_ctx, 7)");
						setState(263);
						((BinaryExprContext)_localctx).op = _input.LT(1);
						_la = _input.LA(1);
						if ( !(_la==PLUS || _la==MINUS) ) {
							((BinaryExprContext)_localctx).op = (Token)_errHandler.recoverInline(this);
						}
						else {
							if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
							_errHandler.reportMatch(this);
							consume();
						}
						setState(264);
						((BinaryExprContext)_localctx).right = expression(8);
						}
						break;
					case 3:
						{
						_localctx = new BinaryExprContext(new ExpressionContext(_parentctx, _parentState));
						((BinaryExprContext)_localctx).left = _prevctx;
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(265);
						if (!(precpred(_ctx, 6))) throw new FailedPredicateException(this, "precpred(_ctx, 6)");
						setState(266);
						((BinaryExprContext)_localctx).op = _input.LT(1);
						_la = _input.LA(1);
						if ( !((((_la) & ~0x3f) == 0 && ((1L << _la) & 2013265920L) != 0)) ) {
							((BinaryExprContext)_localctx).op = (Token)_errHandler.recoverInline(this);
						}
						else {
							if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
							_errHandler.reportMatch(this);
							consume();
						}
						setState(267);
						((BinaryExprContext)_localctx).right = expression(7);
						}
						break;
					case 4:
						{
						_localctx = new BinaryExprContext(new ExpressionContext(_parentctx, _parentState));
						((BinaryExprContext)_localctx).left = _prevctx;
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(268);
						if (!(precpred(_ctx, 5))) throw new FailedPredicateException(this, "precpred(_ctx, 5)");
						setState(269);
						((BinaryExprContext)_localctx).op = _input.LT(1);
						_la = _input.LA(1);
						if ( !(_la==EQ || _la==NE) ) {
							((BinaryExprContext)_localctx).op = (Token)_errHandler.recoverInline(this);
						}
						else {
							if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
							_errHandler.reportMatch(this);
							consume();
						}
						setState(270);
						((BinaryExprContext)_localctx).right = expression(6);
						}
						break;
					case 5:
						{
						_localctx = new BinaryExprContext(new ExpressionContext(_parentctx, _parentState));
						((BinaryExprContext)_localctx).left = _prevctx;
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(271);
						if (!(precpred(_ctx, 4))) throw new FailedPredicateException(this, "precpred(_ctx, 4)");
						setState(272);
						((BinaryExprContext)_localctx).op = match(AND);
						setState(273);
						((BinaryExprContext)_localctx).right = expression(5);
						}
						break;
					case 6:
						{
						_localctx = new BinaryExprContext(new ExpressionContext(_parentctx, _parentState));
						((BinaryExprContext)_localctx).left = _prevctx;
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(274);
						if (!(precpred(_ctx, 3))) throw new FailedPredicateException(this, "precpred(_ctx, 3)");
						setState(275);
						((BinaryExprContext)_localctx).op = match(OR);
						setState(276);
						((BinaryExprContext)_localctx).right = expression(4);
						}
						break;
					case 7:
						{
						_localctx = new StructAccessExprContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(277);
						if (!(precpred(_ctx, 1))) throw new FailedPredicateException(this, "precpred(_ctx, 1)");
						setState(278);
						match(DOT);
						setState(279);
						match(ID);
						}
						break;
					}
					} 
				}
				setState(284);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,17,_ctx);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			unrollRecursionContexts(_parentctx);
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class If_stmtContext extends ParserRuleContext {
		public If_stmtContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_if_stmt; }
	 
		public If_stmtContext() { }
		public void copyFrom(If_stmtContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class IfStmtContext extends If_stmtContext {
		public List<If_chainContext> if_chain() {
			return getRuleContexts(If_chainContext.class);
		}
		public If_chainContext if_chain(int i) {
			return getRuleContext(If_chainContext.class,i);
		}
		public List<TerminalNode> ELSE_KW() { return getTokens(VLangGrammar.ELSE_KW); }
		public TerminalNode ELSE_KW(int i) {
			return getToken(VLangGrammar.ELSE_KW, i);
		}
		public Else_stmtContext else_stmt() {
			return getRuleContext(Else_stmtContext.class,0);
		}
		public IfStmtContext(If_stmtContext ctx) { copyFrom(ctx); }
	}

	public final If_stmtContext if_stmt() throws RecognitionException {
		If_stmtContext _localctx = new If_stmtContext(_ctx, getState());
		enterRule(_localctx, 36, RULE_if_stmt);
		int _la;
		try {
			int _alt;
			_localctx = new IfStmtContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(285);
			if_chain();
			setState(290);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,18,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(286);
					match(ELSE_KW);
					setState(287);
					if_chain();
					}
					} 
				}
				setState(292);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,18,_ctx);
			}
			setState(294);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==ELSE_KW) {
				{
				setState(293);
				else_stmt();
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class If_chainContext extends ParserRuleContext {
		public If_chainContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_if_chain; }
	 
		public If_chainContext() { }
		public void copyFrom(If_chainContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class IfChainContext extends If_chainContext {
		public TerminalNode IF_KW() { return getToken(VLangGrammar.IF_KW, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode LBRACE() { return getToken(VLangGrammar.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(VLangGrammar.RBRACE, 0); }
		public List<StmtContext> stmt() {
			return getRuleContexts(StmtContext.class);
		}
		public StmtContext stmt(int i) {
			return getRuleContext(StmtContext.class,i);
		}
		public IfChainContext(If_chainContext ctx) { copyFrom(ctx); }
	}

	public final If_chainContext if_chain() throws RecognitionException {
		If_chainContext _localctx = new If_chainContext(_ctx, getState());
		enterRule(_localctx, 38, RULE_if_chain);
		int _la;
		try {
			_localctx = new IfChainContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(296);
			match(IF_KW);
			setState(297);
			expression(0);
			setState(298);
			match(LBRACE);
			setState(302);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 563018672928350L) != 0)) {
				{
				{
				setState(299);
				stmt();
				}
				}
				setState(304);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(305);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Else_stmtContext extends ParserRuleContext {
		public Else_stmtContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_else_stmt; }
	 
		public Else_stmtContext() { }
		public void copyFrom(Else_stmtContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class ElseStmtContext extends Else_stmtContext {
		public TerminalNode ELSE_KW() { return getToken(VLangGrammar.ELSE_KW, 0); }
		public TerminalNode LBRACE() { return getToken(VLangGrammar.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(VLangGrammar.RBRACE, 0); }
		public List<StmtContext> stmt() {
			return getRuleContexts(StmtContext.class);
		}
		public StmtContext stmt(int i) {
			return getRuleContext(StmtContext.class,i);
		}
		public ElseStmtContext(Else_stmtContext ctx) { copyFrom(ctx); }
	}

	public final Else_stmtContext else_stmt() throws RecognitionException {
		Else_stmtContext _localctx = new Else_stmtContext(_ctx, getState());
		enterRule(_localctx, 40, RULE_else_stmt);
		int _la;
		try {
			_localctx = new ElseStmtContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(307);
			match(ELSE_KW);
			setState(308);
			match(LBRACE);
			setState(312);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 563018672928350L) != 0)) {
				{
				{
				setState(309);
				stmt();
				}
				}
				setState(314);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(315);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Switch_stmtContext extends ParserRuleContext {
		public Switch_stmtContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_switch_stmt; }
	 
		public Switch_stmtContext() { }
		public void copyFrom(Switch_stmtContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class SwitchStmtContext extends Switch_stmtContext {
		public TerminalNode SWITCH_KW() { return getToken(VLangGrammar.SWITCH_KW, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode LBRACE() { return getToken(VLangGrammar.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(VLangGrammar.RBRACE, 0); }
		public List<Switch_caseContext> switch_case() {
			return getRuleContexts(Switch_caseContext.class);
		}
		public Switch_caseContext switch_case(int i) {
			return getRuleContext(Switch_caseContext.class,i);
		}
		public Default_caseContext default_case() {
			return getRuleContext(Default_caseContext.class,0);
		}
		public SwitchStmtContext(Switch_stmtContext ctx) { copyFrom(ctx); }
	}

	public final Switch_stmtContext switch_stmt() throws RecognitionException {
		Switch_stmtContext _localctx = new Switch_stmtContext(_ctx, getState());
		enterRule(_localctx, 42, RULE_switch_stmt);
		int _la;
		try {
			_localctx = new SwitchStmtContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(317);
			match(SWITCH_KW);
			setState(318);
			expression(0);
			setState(319);
			match(LBRACE);
			setState(323);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==CASE_KW) {
				{
				{
				setState(320);
				switch_case();
				}
				}
				setState(325);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(327);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==DEFAULT_KW) {
				{
				setState(326);
				default_case();
				}
			}

			setState(329);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Switch_caseContext extends ParserRuleContext {
		public Switch_caseContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_switch_case; }
	 
		public Switch_caseContext() { }
		public void copyFrom(Switch_caseContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class SwitchCaseContext extends Switch_caseContext {
		public TerminalNode CASE_KW() { return getToken(VLangGrammar.CASE_KW, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode COLON() { return getToken(VLangGrammar.COLON, 0); }
		public List<StmtContext> stmt() {
			return getRuleContexts(StmtContext.class);
		}
		public StmtContext stmt(int i) {
			return getRuleContext(StmtContext.class,i);
		}
		public SwitchCaseContext(Switch_caseContext ctx) { copyFrom(ctx); }
	}

	public final Switch_caseContext switch_case() throws RecognitionException {
		Switch_caseContext _localctx = new Switch_caseContext(_ctx, getState());
		enterRule(_localctx, 44, RULE_switch_case);
		int _la;
		try {
			_localctx = new SwitchCaseContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(331);
			match(CASE_KW);
			setState(332);
			expression(0);
			setState(333);
			match(COLON);
			setState(337);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 563018672928350L) != 0)) {
				{
				{
				setState(334);
				stmt();
				}
				}
				setState(339);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Default_caseContext extends ParserRuleContext {
		public Default_caseContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_default_case; }
	 
		public Default_caseContext() { }
		public void copyFrom(Default_caseContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class DefaultCaseContext extends Default_caseContext {
		public TerminalNode DEFAULT_KW() { return getToken(VLangGrammar.DEFAULT_KW, 0); }
		public TerminalNode COLON() { return getToken(VLangGrammar.COLON, 0); }
		public List<StmtContext> stmt() {
			return getRuleContexts(StmtContext.class);
		}
		public StmtContext stmt(int i) {
			return getRuleContext(StmtContext.class,i);
		}
		public DefaultCaseContext(Default_caseContext ctx) { copyFrom(ctx); }
	}

	public final Default_caseContext default_case() throws RecognitionException {
		Default_caseContext _localctx = new Default_caseContext(_ctx, getState());
		enterRule(_localctx, 46, RULE_default_case);
		int _la;
		try {
			_localctx = new DefaultCaseContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(340);
			match(DEFAULT_KW);
			setState(341);
			match(COLON);
			setState(345);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 563018672928350L) != 0)) {
				{
				{
				setState(342);
				stmt();
				}
				}
				setState(347);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class While_stmtContext extends ParserRuleContext {
		public While_stmtContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_while_stmt; }
	 
		public While_stmtContext() { }
		public void copyFrom(While_stmtContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class WhileStmtContext extends While_stmtContext {
		public TerminalNode WHILE_KW() { return getToken(VLangGrammar.WHILE_KW, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode LBRACE() { return getToken(VLangGrammar.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(VLangGrammar.RBRACE, 0); }
		public List<StmtContext> stmt() {
			return getRuleContexts(StmtContext.class);
		}
		public StmtContext stmt(int i) {
			return getRuleContext(StmtContext.class,i);
		}
		public WhileStmtContext(While_stmtContext ctx) { copyFrom(ctx); }
	}

	public final While_stmtContext while_stmt() throws RecognitionException {
		While_stmtContext _localctx = new While_stmtContext(_ctx, getState());
		enterRule(_localctx, 48, RULE_while_stmt);
		int _la;
		try {
			_localctx = new WhileStmtContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(348);
			match(WHILE_KW);
			setState(349);
			expression(0);
			setState(350);
			match(LBRACE);
			setState(354);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 563018672928350L) != 0)) {
				{
				{
				setState(351);
				stmt();
				}
				}
				setState(356);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(357);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class For_stmtContext extends ParserRuleContext {
		public For_stmtContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_for_stmt; }
	 
		public For_stmtContext() { }
		public void copyFrom(For_stmtContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class ForStmtCondContext extends For_stmtContext {
		public TerminalNode FOR_KW() { return getToken(VLangGrammar.FOR_KW, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode LBRACE() { return getToken(VLangGrammar.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(VLangGrammar.RBRACE, 0); }
		public List<StmtContext> stmt() {
			return getRuleContexts(StmtContext.class);
		}
		public StmtContext stmt(int i) {
			return getRuleContext(StmtContext.class,i);
		}
		public ForStmtCondContext(For_stmtContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class ForAssCondContext extends For_stmtContext {
		public TerminalNode FOR_KW() { return getToken(VLangGrammar.FOR_KW, 0); }
		public Assign_stmtContext assign_stmt() {
			return getRuleContext(Assign_stmtContext.class,0);
		}
		public List<TerminalNode> SEMI() { return getTokens(VLangGrammar.SEMI); }
		public TerminalNode SEMI(int i) {
			return getToken(VLangGrammar.SEMI, i);
		}
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode LBRACE() { return getToken(VLangGrammar.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(VLangGrammar.RBRACE, 0); }
		public List<StmtContext> stmt() {
			return getRuleContexts(StmtContext.class);
		}
		public StmtContext stmt(int i) {
			return getRuleContext(StmtContext.class,i);
		}
		public ForAssCondContext(For_stmtContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class ForStmtContext extends For_stmtContext {
		public TerminalNode FOR_KW() { return getToken(VLangGrammar.FOR_KW, 0); }
		public List<TerminalNode> ID() { return getTokens(VLangGrammar.ID); }
		public TerminalNode ID(int i) {
			return getToken(VLangGrammar.ID, i);
		}
		public TerminalNode COMMA() { return getToken(VLangGrammar.COMMA, 0); }
		public TerminalNode IN_KW() { return getToken(VLangGrammar.IN_KW, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode LBRACE() { return getToken(VLangGrammar.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(VLangGrammar.RBRACE, 0); }
		public List<StmtContext> stmt() {
			return getRuleContexts(StmtContext.class);
		}
		public StmtContext stmt(int i) {
			return getRuleContext(StmtContext.class,i);
		}
		public ForStmtContext(For_stmtContext ctx) { copyFrom(ctx); }
	}

	public final For_stmtContext for_stmt() throws RecognitionException {
		For_stmtContext _localctx = new For_stmtContext(_ctx, getState());
		enterRule(_localctx, 50, RULE_for_stmt);
		int _la;
		try {
			setState(400);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,30,_ctx) ) {
			case 1:
				_localctx = new ForStmtCondContext(_localctx);
				enterOuterAlt(_localctx, 1);
				{
				setState(359);
				match(FOR_KW);
				setState(360);
				expression(0);
				setState(361);
				match(LBRACE);
				setState(365);
				_errHandler.sync(this);
				_la = _input.LA(1);
				while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 563018672928350L) != 0)) {
					{
					{
					setState(362);
					stmt();
					}
					}
					setState(367);
					_errHandler.sync(this);
					_la = _input.LA(1);
				}
				setState(368);
				match(RBRACE);
				}
				break;
			case 2:
				_localctx = new ForAssCondContext(_localctx);
				enterOuterAlt(_localctx, 2);
				{
				setState(370);
				match(FOR_KW);
				setState(371);
				assign_stmt();
				setState(372);
				match(SEMI);
				setState(373);
				expression(0);
				setState(374);
				match(SEMI);
				setState(375);
				expression(0);
				setState(376);
				match(LBRACE);
				setState(380);
				_errHandler.sync(this);
				_la = _input.LA(1);
				while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 563018672928350L) != 0)) {
					{
					{
					setState(377);
					stmt();
					}
					}
					setState(382);
					_errHandler.sync(this);
					_la = _input.LA(1);
				}
				setState(383);
				match(RBRACE);
				}
				break;
			case 3:
				_localctx = new ForStmtContext(_localctx);
				enterOuterAlt(_localctx, 3);
				{
				setState(385);
				match(FOR_KW);
				setState(386);
				match(ID);
				setState(387);
				match(COMMA);
				setState(388);
				match(ID);
				setState(389);
				match(IN_KW);
				setState(390);
				expression(0);
				setState(391);
				match(LBRACE);
				setState(395);
				_errHandler.sync(this);
				_la = _input.LA(1);
				while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 563018672928350L) != 0)) {
					{
					{
					setState(392);
					stmt();
					}
					}
					setState(397);
					_errHandler.sync(this);
					_la = _input.LA(1);
				}
				setState(398);
				match(RBRACE);
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class RangeContext extends ParserRuleContext {
		public RangeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_range; }
	 
		public RangeContext() { }
		public void copyFrom(RangeContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class NumericRangeContext extends RangeContext {
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public List<TerminalNode> DOT() { return getTokens(VLangGrammar.DOT); }
		public TerminalNode DOT(int i) {
			return getToken(VLangGrammar.DOT, i);
		}
		public NumericRangeContext(RangeContext ctx) { copyFrom(ctx); }
	}

	public final RangeContext range() throws RecognitionException {
		RangeContext _localctx = new RangeContext(_ctx, getState());
		enterRule(_localctx, 52, RULE_range);
		try {
			_localctx = new NumericRangeContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(402);
			expression(0);
			setState(403);
			match(DOT);
			setState(404);
			match(DOT);
			setState(405);
			match(DOT);
			setState(406);
			expression(0);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Transfer_stmtContext extends ParserRuleContext {
		public Transfer_stmtContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_transfer_stmt; }
	 
		public Transfer_stmtContext() { }
		public void copyFrom(Transfer_stmtContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class ContinueStmtContext extends Transfer_stmtContext {
		public TerminalNode CONTINUE_KW() { return getToken(VLangGrammar.CONTINUE_KW, 0); }
		public ContinueStmtContext(Transfer_stmtContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class BreakStmtContext extends Transfer_stmtContext {
		public TerminalNode BREAK_KW() { return getToken(VLangGrammar.BREAK_KW, 0); }
		public BreakStmtContext(Transfer_stmtContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class ReturnStmtContext extends Transfer_stmtContext {
		public TerminalNode RETURN_KW() { return getToken(VLangGrammar.RETURN_KW, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public ReturnStmtContext(Transfer_stmtContext ctx) { copyFrom(ctx); }
	}

	public final Transfer_stmtContext transfer_stmt() throws RecognitionException {
		Transfer_stmtContext _localctx = new Transfer_stmtContext(_ctx, getState());
		enterRule(_localctx, 54, RULE_transfer_stmt);
		try {
			setState(414);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case RETURN_KW:
				_localctx = new ReturnStmtContext(_localctx);
				enterOuterAlt(_localctx, 1);
				{
				setState(408);
				match(RETURN_KW);
				setState(410);
				_errHandler.sync(this);
				switch ( getInterpreter().adaptivePredict(_input,31,_ctx) ) {
				case 1:
					{
					setState(409);
					expression(0);
					}
					break;
				}
				}
				break;
			case BREAK_KW:
				_localctx = new BreakStmtContext(_localctx);
				enterOuterAlt(_localctx, 2);
				{
				setState(412);
				match(BREAK_KW);
				}
				break;
			case CONTINUE_KW:
				_localctx = new ContinueStmtContext(_localctx);
				enterOuterAlt(_localctx, 3);
				{
				setState(413);
				match(CONTINUE_KW);
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Func_callContext extends ParserRuleContext {
		public Func_callContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_func_call; }
	 
		public Func_callContext() { }
		public void copyFrom(Func_callContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class FuncCallContext extends Func_callContext {
		public Id_patternContext id_pattern() {
			return getRuleContext(Id_patternContext.class,0);
		}
		public TerminalNode LPAREN() { return getToken(VLangGrammar.LPAREN, 0); }
		public TerminalNode RPAREN() { return getToken(VLangGrammar.RPAREN, 0); }
		public Arg_listContext arg_list() {
			return getRuleContext(Arg_listContext.class,0);
		}
		public FuncCallContext(Func_callContext ctx) { copyFrom(ctx); }
	}

	public final Func_callContext func_call() throws RecognitionException {
		Func_callContext _localctx = new Func_callContext(_ctx, getState());
		enterRule(_localctx, 56, RULE_func_call);
		int _la;
		try {
			_localctx = new FuncCallContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(416);
			id_pattern();
			setState(417);
			match(LPAREN);
			setState(419);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if ((((_la) & ~0x3f) == 0 && ((1L << _la) & 1108677088247808L) != 0)) {
				{
				setState(418);
				arg_list();
				}
			}

			setState(421);
			match(RPAREN);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Block_indContext extends ParserRuleContext {
		public Block_indContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_block_ind; }
	 
		public Block_indContext() { }
		public void copyFrom(Block_indContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class BlockIndContext extends Block_indContext {
		public TerminalNode LBRACE() { return getToken(VLangGrammar.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(VLangGrammar.RBRACE, 0); }
		public List<StmtContext> stmt() {
			return getRuleContexts(StmtContext.class);
		}
		public StmtContext stmt(int i) {
			return getRuleContext(StmtContext.class,i);
		}
		public BlockIndContext(Block_indContext ctx) { copyFrom(ctx); }
	}

	public final Block_indContext block_ind() throws RecognitionException {
		Block_indContext _localctx = new Block_indContext(_ctx, getState());
		enterRule(_localctx, 58, RULE_block_ind);
		int _la;
		try {
			_localctx = new BlockIndContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(423);
			match(LBRACE);
			setState(427);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 563018672928350L) != 0)) {
				{
				{
				setState(424);
				stmt();
				}
				}
				setState(429);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(430);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Arg_listContext extends ParserRuleContext {
		public Arg_listContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_arg_list; }
	 
		public Arg_listContext() { }
		public void copyFrom(Arg_listContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class ArgListContext extends Arg_listContext {
		public List<Func_argContext> func_arg() {
			return getRuleContexts(Func_argContext.class);
		}
		public Func_argContext func_arg(int i) {
			return getRuleContext(Func_argContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(VLangGrammar.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(VLangGrammar.COMMA, i);
		}
		public ArgListContext(Arg_listContext ctx) { copyFrom(ctx); }
	}

	public final Arg_listContext arg_list() throws RecognitionException {
		Arg_listContext _localctx = new Arg_listContext(_ctx, getState());
		enterRule(_localctx, 60, RULE_arg_list);
		int _la;
		try {
			_localctx = new ArgListContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(432);
			func_arg();
			setState(437);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(433);
				match(COMMA);
				setState(434);
				func_arg();
				}
				}
				setState(439);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Func_argContext extends ParserRuleContext {
		public Func_argContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_func_arg; }
	 
		public Func_argContext() { }
		public void copyFrom(Func_argContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class FuncArgContext extends Func_argContext {
		public Id_patternContext id_pattern() {
			return getRuleContext(Id_patternContext.class,0);
		}
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public FuncArgContext(Func_argContext ctx) { copyFrom(ctx); }
	}

	public final Func_argContext func_arg() throws RecognitionException {
		Func_argContext _localctx = new Func_argContext(_ctx, getState());
		enterRule(_localctx, 62, RULE_func_arg);
		try {
			_localctx = new FuncArgContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(441);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,36,_ctx) ) {
			case 1:
				{
				setState(440);
				match(ID);
				}
				break;
			}
			setState(445);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,37,_ctx) ) {
			case 1:
				{
				setState(443);
				id_pattern();
				}
				break;
			case 2:
				{
				setState(444);
				expression(0);
				}
				break;
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Func_dclContext extends ParserRuleContext {
		public Func_dclContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_func_dcl; }
	 
		public Func_dclContext() { }
		public void copyFrom(Func_dclContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class FuncDeclContext extends Func_dclContext {
		public TerminalNode FUNC() { return getToken(VLangGrammar.FUNC, 0); }
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public TerminalNode LPAREN() { return getToken(VLangGrammar.LPAREN, 0); }
		public TerminalNode RPAREN() { return getToken(VLangGrammar.RPAREN, 0); }
		public TerminalNode LBRACE() { return getToken(VLangGrammar.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(VLangGrammar.RBRACE, 0); }
		public Param_listContext param_list() {
			return getRuleContext(Param_listContext.class,0);
		}
		public TypeContext type() {
			return getRuleContext(TypeContext.class,0);
		}
		public List<StmtContext> stmt() {
			return getRuleContexts(StmtContext.class);
		}
		public StmtContext stmt(int i) {
			return getRuleContext(StmtContext.class,i);
		}
		public FuncDeclContext(Func_dclContext ctx) { copyFrom(ctx); }
	}

	public final Func_dclContext func_dcl() throws RecognitionException {
		Func_dclContext _localctx = new Func_dclContext(_ctx, getState());
		enterRule(_localctx, 64, RULE_func_dcl);
		int _la;
		try {
			_localctx = new FuncDeclContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(447);
			match(FUNC);
			setState(448);
			match(ID);
			setState(449);
			match(LPAREN);
			setState(451);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==ID) {
				{
				setState(450);
				param_list();
				}
			}

			setState(453);
			match(RPAREN);
			setState(455);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==LBRACK || _la==ID) {
				{
				setState(454);
				type();
				}
			}

			setState(457);
			match(LBRACE);
			setState(461);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 563018672928350L) != 0)) {
				{
				{
				setState(458);
				stmt();
				}
				}
				setState(463);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(464);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Param_listContext extends ParserRuleContext {
		public Param_listContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_param_list; }
	 
		public Param_listContext() { }
		public void copyFrom(Param_listContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class ParamListContext extends Param_listContext {
		public List<Func_paramContext> func_param() {
			return getRuleContexts(Func_paramContext.class);
		}
		public Func_paramContext func_param(int i) {
			return getRuleContext(Func_paramContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(VLangGrammar.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(VLangGrammar.COMMA, i);
		}
		public ParamListContext(Param_listContext ctx) { copyFrom(ctx); }
	}

	public final Param_listContext param_list() throws RecognitionException {
		Param_listContext _localctx = new Param_listContext(_ctx, getState());
		enterRule(_localctx, 66, RULE_param_list);
		int _la;
		try {
			_localctx = new ParamListContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(466);
			func_param();
			setState(471);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(467);
				match(COMMA);
				setState(468);
				func_param();
				}
				}
				setState(473);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Func_paramContext extends ParserRuleContext {
		public Func_paramContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_func_param; }
	 
		public Func_paramContext() { }
		public void copyFrom(Func_paramContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class FuncParamContext extends Func_paramContext {
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public TypeContext type() {
			return getRuleContext(TypeContext.class,0);
		}
		public FuncParamContext(Func_paramContext ctx) { copyFrom(ctx); }
	}

	public final Func_paramContext func_param() throws RecognitionException {
		Func_paramContext _localctx = new Func_paramContext(_ctx, getState());
		enterRule(_localctx, 68, RULE_func_param);
		try {
			_localctx = new FuncParamContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(474);
			match(ID);
			setState(475);
			type();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Strct_dclContext extends ParserRuleContext {
		public Strct_dclContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_strct_dcl; }
	 
		public Strct_dclContext() { }
		public void copyFrom(Strct_dclContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class StructDeclContext extends Strct_dclContext {
		public TerminalNode STR() { return getToken(VLangGrammar.STR, 0); }
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public TerminalNode LBRACE() { return getToken(VLangGrammar.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(VLangGrammar.RBRACE, 0); }
		public List<Struct_propContext> struct_prop() {
			return getRuleContexts(Struct_propContext.class);
		}
		public Struct_propContext struct_prop(int i) {
			return getRuleContext(Struct_propContext.class,i);
		}
		public StructDeclContext(Strct_dclContext ctx) { copyFrom(ctx); }
	}

	public final Strct_dclContext strct_dcl() throws RecognitionException {
		Strct_dclContext _localctx = new Strct_dclContext(_ctx, getState());
		enterRule(_localctx, 70, RULE_strct_dcl);
		int _la;
		try {
			_localctx = new StructDeclContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(477);
			match(STR);
			setState(478);
			match(ID);
			setState(479);
			match(LBRACE);
			setState(481); 
			_errHandler.sync(this);
			_la = _input.LA(1);
			do {
				{
				{
				setState(480);
				struct_prop();
				}
				}
				setState(483); 
				_errHandler.sync(this);
				_la = _input.LA(1);
			} while ( _la==LBRACK || _la==ID );
			setState(485);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Struct_propContext extends ParserRuleContext {
		public Struct_propContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_struct_prop; }
	 
		public Struct_propContext() { }
		public void copyFrom(Struct_propContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class StructAttrContext extends Struct_propContext {
		public TypeContext type() {
			return getRuleContext(TypeContext.class,0);
		}
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public TerminalNode SEMI() { return getToken(VLangGrammar.SEMI, 0); }
		public StructAttrContext(Struct_propContext ctx) { copyFrom(ctx); }
	}

	public final Struct_propContext struct_prop() throws RecognitionException {
		Struct_propContext _localctx = new Struct_propContext(_ctx, getState());
		enterRule(_localctx, 72, RULE_struct_prop);
		try {
			_localctx = new StructAttrContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(487);
			type();
			setState(488);
			match(ID);
			setState(489);
			match(SEMI);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Struct_param_listContext extends ParserRuleContext {
		public List<Struct_paramContext> struct_param() {
			return getRuleContexts(Struct_paramContext.class);
		}
		public Struct_paramContext struct_param(int i) {
			return getRuleContext(Struct_paramContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(VLangGrammar.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(VLangGrammar.COMMA, i);
		}
		public Struct_param_listContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_struct_param_list; }
	}

	public final Struct_param_listContext struct_param_list() throws RecognitionException {
		Struct_param_listContext _localctx = new Struct_param_listContext(_ctx, getState());
		enterRule(_localctx, 74, RULE_struct_param_list);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(491);
			struct_param();
			setState(496);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,43,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(492);
					match(COMMA);
					setState(493);
					struct_param();
					}
					} 
				}
				setState(498);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,43,_ctx);
			}
			setState(500);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==COMMA) {
				{
				setState(499);
				match(COMMA);
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class Struct_paramContext extends ParserRuleContext {
		public TerminalNode ID() { return getToken(VLangGrammar.ID, 0); }
		public TerminalNode COLON() { return getToken(VLangGrammar.COLON, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public Struct_paramContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_struct_param; }
	}

	public final Struct_paramContext struct_param() throws RecognitionException {
		Struct_paramContext _localctx = new Struct_paramContext(_ctx, getState());
		enterRule(_localctx, 76, RULE_struct_param);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(502);
			match(ID);
			setState(503);
			match(COLON);
			setState(504);
			expression(0);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public boolean sempred(RuleContext _localctx, int ruleIndex, int predIndex) {
		switch (ruleIndex) {
		case 17:
			return expression_sempred((ExpressionContext)_localctx, predIndex);
		}
		return true;
	}
	private boolean expression_sempred(ExpressionContext _localctx, int predIndex) {
		switch (predIndex) {
		case 0:
			return precpred(_ctx, 8);
		case 1:
			return precpred(_ctx, 7);
		case 2:
			return precpred(_ctx, 6);
		case 3:
			return precpred(_ctx, 5);
		case 4:
			return precpred(_ctx, 4);
		case 5:
			return precpred(_ctx, 3);
		case 6:
			return precpred(_ctx, 1);
		}
		return true;
	}

	public static final String _serializedATN =
		"\u0004\u00014\u01fb\u0002\u0000\u0007\u0000\u0002\u0001\u0007\u0001\u0002"+
		"\u0002\u0007\u0002\u0002\u0003\u0007\u0003\u0002\u0004\u0007\u0004\u0002"+
		"\u0005\u0007\u0005\u0002\u0006\u0007\u0006\u0002\u0007\u0007\u0007\u0002"+
		"\b\u0007\b\u0002\t\u0007\t\u0002\n\u0007\n\u0002\u000b\u0007\u000b\u0002"+
		"\f\u0007\f\u0002\r\u0007\r\u0002\u000e\u0007\u000e\u0002\u000f\u0007\u000f"+
		"\u0002\u0010\u0007\u0010\u0002\u0011\u0007\u0011\u0002\u0012\u0007\u0012"+
		"\u0002\u0013\u0007\u0013\u0002\u0014\u0007\u0014\u0002\u0015\u0007\u0015"+
		"\u0002\u0016\u0007\u0016\u0002\u0017\u0007\u0017\u0002\u0018\u0007\u0018"+
		"\u0002\u0019\u0007\u0019\u0002\u001a\u0007\u001a\u0002\u001b\u0007\u001b"+
		"\u0002\u001c\u0007\u001c\u0002\u001d\u0007\u001d\u0002\u001e\u0007\u001e"+
		"\u0002\u001f\u0007\u001f\u0002 \u0007 \u0002!\u0007!\u0002\"\u0007\"\u0002"+
		"#\u0007#\u0002$\u0007$\u0002%\u0007%\u0002&\u0007&\u0001\u0000\u0005\u0000"+
		"P\b\u0000\n\u0000\f\u0000S\t\u0000\u0001\u0000\u0003\u0000V\b\u0000\u0001"+
		"\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001"+
		"\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0003"+
		"\u0001d\b\u0001\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001"+
		"\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001"+
		"\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001"+
		"\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001"+
		"\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001"+
		"\u0002\u0001\u0002\u0003\u0002\u0084\b\u0002\u0001\u0003\u0001\u0003\u0001"+
		"\u0004\u0001\u0004\u0001\u0004\u0001\u0004\u0005\u0004\u008c\b\u0004\n"+
		"\u0004\f\u0004\u008f\t\u0004\u0003\u0004\u0091\b\u0004\u0001\u0004\u0001"+
		"\u0004\u0001\u0005\u0001\u0005\u0001\u0005\u0001\u0005\u0001\u0005\u0004"+
		"\u0005\u009a\b\u0005\u000b\u0005\f\u0005\u009b\u0001\u0006\u0001\u0006"+
		"\u0001\u0006\u0001\u0006\u0001\u0007\u0001\u0007\u0001\u0007\u0001\u0007"+
		"\u0001\b\u0001\b\u0003\b\u00a8\b\b\u0001\b\u0001\b\u0001\b\u0001\b\u0001"+
		"\b\u0001\b\u0001\b\u0001\b\u0001\b\u0001\b\u0001\t\u0001\t\u0001\t\u0001"+
		"\t\u0001\n\u0001\n\u0001\n\u0001\n\u0001\n\u0001\n\u0003\n\u00be\b\n\u0001"+
		"\u000b\u0001\u000b\u0001\u000b\u0001\u000b\u0001\f\u0001\f\u0001\f\u0003"+
		"\f\u00c7\b\f\u0001\r\u0001\r\u0001\r\u0001\r\u0001\r\u0001\r\u0001\r\u0001"+
		"\r\u0001\r\u0001\r\u0001\r\u0001\r\u0003\r\u00d5\b\r\u0001\u000e\u0001"+
		"\u000e\u0001\u000e\u0005\u000e\u00da\b\u000e\n\u000e\f\u000e\u00dd\t\u000e"+
		"\u0001\u000f\u0001\u000f\u0001\u000f\u0001\u000f\u0001\u000f\u0003\u000f"+
		"\u00e4\b\u000f\u0001\u0010\u0001\u0010\u0001\u0010\u0001\u0010\u0003\u0010"+
		"\u00ea\b\u0010\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011"+
		"\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011"+
		"\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011"+
		"\u0001\u0011\u0001\u0011\u0003\u0011\u00ff\b\u0011\u0001\u0011\u0003\u0011"+
		"\u0102\b\u0011\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011"+
		"\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011"+
		"\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011"+
		"\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011\u0005\u0011\u0119\b\u0011"+
		"\n\u0011\f\u0011\u011c\t\u0011\u0001\u0012\u0001\u0012\u0001\u0012\u0005"+
		"\u0012\u0121\b\u0012\n\u0012\f\u0012\u0124\t\u0012\u0001\u0012\u0003\u0012"+
		"\u0127\b\u0012\u0001\u0013\u0001\u0013\u0001\u0013\u0001\u0013\u0005\u0013"+
		"\u012d\b\u0013\n\u0013\f\u0013\u0130\t\u0013\u0001\u0013\u0001\u0013\u0001"+
		"\u0014\u0001\u0014\u0001\u0014\u0005\u0014\u0137\b\u0014\n\u0014\f\u0014"+
		"\u013a\t\u0014\u0001\u0014\u0001\u0014\u0001\u0015\u0001\u0015\u0001\u0015"+
		"\u0001\u0015\u0005\u0015\u0142\b\u0015\n\u0015\f\u0015\u0145\t\u0015\u0001"+
		"\u0015\u0003\u0015\u0148\b\u0015\u0001\u0015\u0001\u0015\u0001\u0016\u0001"+
		"\u0016\u0001\u0016\u0001\u0016\u0005\u0016\u0150\b\u0016\n\u0016\f\u0016"+
		"\u0153\t\u0016\u0001\u0017\u0001\u0017\u0001\u0017\u0005\u0017\u0158\b"+
		"\u0017\n\u0017\f\u0017\u015b\t\u0017\u0001\u0018\u0001\u0018\u0001\u0018"+
		"\u0001\u0018\u0005\u0018\u0161\b\u0018\n\u0018\f\u0018\u0164\t\u0018\u0001"+
		"\u0018\u0001\u0018\u0001\u0019\u0001\u0019\u0001\u0019\u0001\u0019\u0005"+
		"\u0019\u016c\b\u0019\n\u0019\f\u0019\u016f\t\u0019\u0001\u0019\u0001\u0019"+
		"\u0001\u0019\u0001\u0019\u0001\u0019\u0001\u0019\u0001\u0019\u0001\u0019"+
		"\u0001\u0019\u0001\u0019\u0005\u0019\u017b\b\u0019\n\u0019\f\u0019\u017e"+
		"\t\u0019\u0001\u0019\u0001\u0019\u0001\u0019\u0001\u0019\u0001\u0019\u0001"+
		"\u0019\u0001\u0019\u0001\u0019\u0001\u0019\u0001\u0019\u0005\u0019\u018a"+
		"\b\u0019\n\u0019\f\u0019\u018d\t\u0019\u0001\u0019\u0001\u0019\u0003\u0019"+
		"\u0191\b\u0019\u0001\u001a\u0001\u001a\u0001\u001a\u0001\u001a\u0001\u001a"+
		"\u0001\u001a\u0001\u001b\u0001\u001b\u0003\u001b\u019b\b\u001b\u0001\u001b"+
		"\u0001\u001b\u0003\u001b\u019f\b\u001b\u0001\u001c\u0001\u001c\u0001\u001c"+
		"\u0003\u001c\u01a4\b\u001c\u0001\u001c\u0001\u001c\u0001\u001d\u0001\u001d"+
		"\u0005\u001d\u01aa\b\u001d\n\u001d\f\u001d\u01ad\t\u001d\u0001\u001d\u0001"+
		"\u001d\u0001\u001e\u0001\u001e\u0001\u001e\u0005\u001e\u01b4\b\u001e\n"+
		"\u001e\f\u001e\u01b7\t\u001e\u0001\u001f\u0003\u001f\u01ba\b\u001f\u0001"+
		"\u001f\u0001\u001f\u0003\u001f\u01be\b\u001f\u0001 \u0001 \u0001 \u0001"+
		" \u0003 \u01c4\b \u0001 \u0001 \u0003 \u01c8\b \u0001 \u0001 \u0005 \u01cc"+
		"\b \n \f \u01cf\t \u0001 \u0001 \u0001!\u0001!\u0001!\u0005!\u01d6\b!"+
		"\n!\f!\u01d9\t!\u0001\"\u0001\"\u0001\"\u0001#\u0001#\u0001#\u0001#\u0004"+
		"#\u01e2\b#\u000b#\f#\u01e3\u0001#\u0001#\u0001$\u0001$\u0001$\u0001$\u0001"+
		"%\u0001%\u0001%\u0005%\u01ef\b%\n%\f%\u01f2\t%\u0001%\u0003%\u01f5\b%"+
		"\u0001&\u0001&\u0001&\u0001&\u0001&\u0000\u0001\"\'\u0000\u0002\u0004"+
		"\u0006\b\n\f\u000e\u0010\u0012\u0014\u0016\u0018\u001a\u001c\u001e \""+
		"$&(*,.02468:<>@BDFHJL\u0000\u0007\u0001\u0000\u0017\u0018\u0001\u0000"+
		"\u0016\u0018\u0002\u0000\u0012\u0012!!\u0001\u0000\u0013\u0015\u0001\u0000"+
		"\u0011\u0012\u0001\u0000\u001b\u001e\u0001\u0000\u0019\u001a\u0224\u0000"+
		"Q\u0001\u0000\u0000\u0000\u0002c\u0001\u0000\u0000\u0000\u0004\u0083\u0001"+
		"\u0000\u0000\u0000\u0006\u0085\u0001\u0000\u0000\u0000\b\u0087\u0001\u0000"+
		"\u0000\u0000\n\u0094\u0001\u0000\u0000\u0000\f\u009d\u0001\u0000\u0000"+
		"\u0000\u000e\u00a1\u0001\u0000\u0000\u0000\u0010\u00a7\u0001\u0000\u0000"+
		"\u0000\u0012\u00b3\u0001\u0000\u0000\u0000\u0014\u00bd\u0001\u0000\u0000"+
		"\u0000\u0016\u00bf\u0001\u0000\u0000\u0000\u0018\u00c6\u0001\u0000\u0000"+
		"\u0000\u001a\u00d4\u0001\u0000\u0000\u0000\u001c\u00d6\u0001\u0000\u0000"+
		"\u0000\u001e\u00e3\u0001\u0000\u0000\u0000 \u00e9\u0001\u0000\u0000\u0000"+
		"\"\u0101\u0001\u0000\u0000\u0000$\u011d\u0001\u0000\u0000\u0000&\u0128"+
		"\u0001\u0000\u0000\u0000(\u0133\u0001\u0000\u0000\u0000*\u013d\u0001\u0000"+
		"\u0000\u0000,\u014b\u0001\u0000\u0000\u0000.\u0154\u0001\u0000\u0000\u0000"+
		"0\u015c\u0001\u0000\u0000\u00002\u0190\u0001\u0000\u0000\u00004\u0192"+
		"\u0001\u0000\u0000\u00006\u019e\u0001\u0000\u0000\u00008\u01a0\u0001\u0000"+
		"\u0000\u0000:\u01a7\u0001\u0000\u0000\u0000<\u01b0\u0001\u0000\u0000\u0000"+
		">\u01b9\u0001\u0000\u0000\u0000@\u01bf\u0001\u0000\u0000\u0000B\u01d2"+
		"\u0001\u0000\u0000\u0000D\u01da\u0001\u0000\u0000\u0000F\u01dd\u0001\u0000"+
		"\u0000\u0000H\u01e7\u0001\u0000\u0000\u0000J\u01eb\u0001\u0000\u0000\u0000"+
		"L\u01f6\u0001\u0000\u0000\u0000NP\u0003\u0002\u0001\u0000ON\u0001\u0000"+
		"\u0000\u0000PS\u0001\u0000\u0000\u0000QO\u0001\u0000\u0000\u0000QR\u0001"+
		"\u0000\u0000\u0000RU\u0001\u0000\u0000\u0000SQ\u0001\u0000\u0000\u0000"+
		"TV\u0005\u0000\u0000\u0001UT\u0001\u0000\u0000\u0000UV\u0001\u0000\u0000"+
		"\u0000V\u0001\u0001\u0000\u0000\u0000Wd\u0003\u0004\u0002\u0000Xd\u0003"+
		"\u001a\r\u0000Yd\u0003:\u001d\u0000Zd\u00036\u001b\u0000[d\u0003$\u0012"+
		"\u0000\\d\u0003*\u0015\u0000]d\u00030\u0018\u0000^d\u00032\u0019\u0000"+
		"_d\u00038\u001c\u0000`d\u0003\u000e\u0007\u0000ad\u0003@ \u0000bd\u0003"+
		"F#\u0000cW\u0001\u0000\u0000\u0000cX\u0001\u0000\u0000\u0000cY\u0001\u0000"+
		"\u0000\u0000cZ\u0001\u0000\u0000\u0000c[\u0001\u0000\u0000\u0000c\\\u0001"+
		"\u0000\u0000\u0000c]\u0001\u0000\u0000\u0000c^\u0001\u0000\u0000\u0000"+
		"c_\u0001\u0000\u0000\u0000c`\u0001\u0000\u0000\u0000ca\u0001\u0000\u0000"+
		"\u0000cb\u0001\u0000\u0000\u0000d\u0003\u0001\u0000\u0000\u0000ef\u0003"+
		"\u0006\u0003\u0000fg\u00051\u0000\u0000gh\u0003\u0018\f\u0000hi\u0005"+
		"\u0016\u0000\u0000ij\u0003\"\u0011\u0000j\u0084\u0001\u0000\u0000\u0000"+
		"kl\u0003\u0006\u0003\u0000lm\u00051\u0000\u0000mn\u0005\u0016\u0000\u0000"+
		"no\u0003\"\u0011\u0000o\u0084\u0001\u0000\u0000\u0000pq\u0003\u0006\u0003"+
		"\u0000qr\u00051\u0000\u0000rs\u0003\u0018\f\u0000s\u0084\u0001\u0000\u0000"+
		"\u0000tu\u00051\u0000\u0000uv\u0003\u0018\f\u0000vw\u0005\u0016\u0000"+
		"\u0000wx\u0003\"\u0011\u0000x\u0084\u0001\u0000\u0000\u0000yz\u00051\u0000"+
		"\u0000z{\u0005\u0016\u0000\u0000{|\u0003\u0012\t\u0000|}\u0003\b\u0004"+
		"\u0000}\u0084\u0001\u0000\u0000\u0000~\u007f\u00051\u0000\u0000\u007f"+
		"\u0080\u0005\u0016\u0000\u0000\u0080\u0081\u0003\u0014\n\u0000\u0081\u0082"+
		"\u0003\b\u0004\u0000\u0082\u0084\u0001\u0000\u0000\u0000\u0083e\u0001"+
		"\u0000\u0000\u0000\u0083k\u0001\u0000\u0000\u0000\u0083p\u0001\u0000\u0000"+
		"\u0000\u0083t\u0001\u0000\u0000\u0000\u0083y\u0001\u0000\u0000\u0000\u0083"+
		"~\u0001\u0000\u0000\u0000\u0084\u0005\u0001\u0000\u0000\u0000\u0085\u0086"+
		"\u0005\u0001\u0000\u0000\u0086\u0007\u0001\u0000\u0000\u0000\u0087\u0090"+
		"\u0005$\u0000\u0000\u0088\u008d\u0003\"\u0011\u0000\u0089\u008a\u0005"+
		"+\u0000\u0000\u008a\u008c\u0003\"\u0011\u0000\u008b\u0089\u0001\u0000"+
		"\u0000\u0000\u008c\u008f\u0001\u0000\u0000\u0000\u008d\u008b\u0001\u0000"+
		"\u0000\u0000\u008d\u008e\u0001\u0000\u0000\u0000\u008e\u0091\u0001\u0000"+
		"\u0000\u0000\u008f\u008d\u0001\u0000\u0000\u0000\u0090\u0088\u0001\u0000"+
		"\u0000\u0000\u0090\u0091\u0001\u0000\u0000\u0000\u0091\u0092\u0001\u0000"+
		"\u0000\u0000\u0092\u0093\u0005%\u0000\u0000\u0093\t\u0001\u0000\u0000"+
		"\u0000\u0094\u0099\u0003\u001c\u000e\u0000\u0095\u0096\u0005&\u0000\u0000"+
		"\u0096\u0097\u0003\"\u0011\u0000\u0097\u0098\u0005\'\u0000\u0000\u0098"+
		"\u009a\u0001\u0000\u0000\u0000\u0099\u0095\u0001\u0000\u0000\u0000\u009a"+
		"\u009b\u0001\u0000\u0000\u0000\u009b\u0099\u0001\u0000\u0000\u0000\u009b"+
		"\u009c\u0001\u0000\u0000\u0000\u009c\u000b\u0001\u0000\u0000\u0000\u009d"+
		"\u009e\u0003\n\u0005\u0000\u009e\u009f\u0005*\u0000\u0000\u009f\u00a0"+
		"\u0003\u001c\u000e\u0000\u00a0\r\u0001\u0000\u0000\u0000\u00a1\u00a2\u0003"+
		"\n\u0005\u0000\u00a2\u00a3\u0005*\u0000\u0000\u00a3\u00a4\u00038\u001c"+
		"\u0000\u00a4\u000f\u0001\u0000\u0000\u0000\u00a5\u00a8\u0003\u0012\t\u0000"+
		"\u00a6\u00a8\u0003\u0014\n\u0000\u00a7\u00a5\u0001\u0000\u0000\u0000\u00a7"+
		"\u00a6\u0001\u0000\u0000\u0000\u00a8\u00a9\u0001\u0000\u0000\u0000\u00a9"+
		"\u00aa\u0005\"\u0000\u0000\u00aa\u00ab\u00051\u0000\u0000\u00ab\u00ac"+
		"\u0005)\u0000\u0000\u00ac\u00ad\u0003\"\u0011\u0000\u00ad\u00ae\u0005"+
		"+\u0000\u0000\u00ae\u00af\u00051\u0000\u0000\u00af\u00b0\u0005)\u0000"+
		"\u0000\u00b0\u00b1\u0003\"\u0011\u0000\u00b1\u00b2\u0005#\u0000\u0000"+
		"\u00b2\u0011\u0001\u0000\u0000\u0000\u00b3\u00b4\u0005&\u0000\u0000\u00b4"+
		"\u00b5\u0005\'\u0000\u0000\u00b5\u00b6\u00051\u0000\u0000\u00b6\u0013"+
		"\u0001\u0000\u0000\u0000\u00b7\u00be\u0003\u0016\u000b\u0000\u00b8\u00b9"+
		"\u0005&\u0000\u0000\u00b9\u00ba\u0005\'\u0000\u0000\u00ba\u00bb\u0005"+
		"&\u0000\u0000\u00bb\u00bc\u0005\'\u0000\u0000\u00bc\u00be\u00051\u0000"+
		"\u0000\u00bd\u00b7\u0001\u0000\u0000\u0000\u00bd\u00b8\u0001\u0000\u0000"+
		"\u0000\u00be\u0015\u0001\u0000\u0000\u0000\u00bf\u00c0\u0005&\u0000\u0000"+
		"\u00c0\u00c1\u0005\'\u0000\u0000\u00c1\u00c2\u0003\u0014\n\u0000\u00c2"+
		"\u0017\u0001\u0000\u0000\u0000\u00c3\u00c7\u00051\u0000\u0000\u00c4\u00c7"+
		"\u0003\u0012\t\u0000\u00c5\u00c7\u0003\u0014\n\u0000\u00c6\u00c3\u0001"+
		"\u0000\u0000\u0000\u00c6\u00c4\u0001\u0000\u0000\u0000\u00c6\u00c5\u0001"+
		"\u0000\u0000\u0000\u00c7\u0019\u0001\u0000\u0000\u0000\u00c8\u00c9\u0003"+
		"\u001c\u000e\u0000\u00c9\u00ca\u0005\u0016\u0000\u0000\u00ca\u00cb\u0003"+
		"\"\u0011\u0000\u00cb\u00d5\u0001\u0000\u0000\u0000\u00cc\u00cd\u0003\u001c"+
		"\u000e\u0000\u00cd\u00ce\u0007\u0000\u0000\u0000\u00ce\u00cf\u0003\"\u0011"+
		"\u0000\u00cf\u00d5\u0001\u0000\u0000\u0000\u00d0\u00d1\u0003\n\u0005\u0000"+
		"\u00d1\u00d2\u0007\u0001\u0000\u0000\u00d2\u00d3\u0003\"\u0011\u0000\u00d3"+
		"\u00d5\u0001\u0000\u0000\u0000\u00d4\u00c8\u0001\u0000\u0000\u0000\u00d4"+
		"\u00cc\u0001\u0000\u0000\u0000\u00d4\u00d0\u0001\u0000\u0000\u0000\u00d5"+
		"\u001b\u0001\u0000\u0000\u0000\u00d6\u00db\u00051\u0000\u0000\u00d7\u00d8"+
		"\u0005*\u0000\u0000\u00d8\u00da\u00051\u0000\u0000\u00d9\u00d7\u0001\u0000"+
		"\u0000\u0000\u00da\u00dd\u0001\u0000\u0000\u0000\u00db\u00d9\u0001\u0000"+
		"\u0000\u0000\u00db\u00dc\u0001\u0000\u0000\u0000\u00dc\u001d\u0001\u0000"+
		"\u0000\u0000\u00dd\u00db\u0001\u0000\u0000\u0000\u00de\u00e4\u0005,\u0000"+
		"\u0000\u00df\u00e4\u0005-\u0000\u0000\u00e0\u00e4\u0005.\u0000\u0000\u00e1"+
		"\u00e4\u0005/\u0000\u0000\u00e2\u00e4\u00050\u0000\u0000\u00e3\u00de\u0001"+
		"\u0000\u0000\u0000\u00e3\u00df\u0001\u0000\u0000\u0000\u00e3\u00e0\u0001"+
		"\u0000\u0000\u0000\u00e3\u00e1\u0001\u0000\u0000\u0000\u00e3\u00e2\u0001"+
		"\u0000\u0000\u0000\u00e4\u001f\u0001\u0000\u0000\u0000\u00e5\u00e6\u0005"+
		"1\u0000\u0000\u00e6\u00ea\u0005\u0010\u0000\u0000\u00e7\u00e8\u00051\u0000"+
		"\u0000\u00e8\u00ea\u0005\u000f\u0000\u0000\u00e9\u00e5\u0001\u0000\u0000"+
		"\u0000\u00e9\u00e7\u0001\u0000\u0000\u0000\u00ea!\u0001\u0000\u0000\u0000"+
		"\u00eb\u00ec\u0006\u0011\uffff\uffff\u0000\u00ec\u00ed\u0005\"\u0000\u0000"+
		"\u00ed\u00ee\u0003\"\u0011\u0000\u00ee\u00ef\u0005#\u0000\u0000\u00ef"+
		"\u0102\u0001\u0000\u0000\u0000\u00f0\u0102\u00038\u001c\u0000\u00f1\u0102"+
		"\u0003\u001c\u000e\u0000\u00f2\u0102\u0003\n\u0005\u0000\u00f3\u0102\u0003"+
		"\f\u0006\u0000\u00f4\u0102\u0003\u000e\u0007\u0000\u00f5\u0102\u0003\u001e"+
		"\u000f\u0000\u00f6\u0102\u0003\b\u0004\u0000\u00f7\u0102\u0003\u0010\b"+
		"\u0000\u00f8\u0102\u0003 \u0010\u0000\u00f9\u00fa\u0007\u0002\u0000\u0000"+
		"\u00fa\u0102\u0003\"\u0011\t\u00fb\u00fc\u00051\u0000\u0000\u00fc\u00fe"+
		"\u0005$\u0000\u0000\u00fd\u00ff\u0003J%\u0000\u00fe\u00fd\u0001\u0000"+
		"\u0000\u0000\u00fe\u00ff\u0001\u0000\u0000\u0000\u00ff\u0100\u0001\u0000"+
		"\u0000\u0000\u0100\u0102\u0005%\u0000\u0000\u0101\u00eb\u0001\u0000\u0000"+
		"\u0000\u0101\u00f0\u0001\u0000\u0000\u0000\u0101\u00f1\u0001\u0000\u0000"+
		"\u0000\u0101\u00f2\u0001\u0000\u0000\u0000\u0101\u00f3\u0001\u0000\u0000"+
		"\u0000\u0101\u00f4\u0001\u0000\u0000\u0000\u0101\u00f5\u0001\u0000\u0000"+
		"\u0000\u0101\u00f6\u0001\u0000\u0000\u0000\u0101\u00f7\u0001\u0000\u0000"+
		"\u0000\u0101\u00f8\u0001\u0000\u0000\u0000\u0101\u00f9\u0001\u0000\u0000"+
		"\u0000\u0101\u00fb\u0001\u0000\u0000\u0000\u0102\u011a\u0001\u0000\u0000"+
		"\u0000\u0103\u0104\n\b\u0000\u0000\u0104\u0105\u0007\u0003\u0000\u0000"+
		"\u0105\u0119\u0003\"\u0011\t\u0106\u0107\n\u0007\u0000\u0000\u0107\u0108"+
		"\u0007\u0004\u0000\u0000\u0108\u0119\u0003\"\u0011\b\u0109\u010a\n\u0006"+
		"\u0000\u0000\u010a\u010b\u0007\u0005\u0000\u0000\u010b\u0119\u0003\"\u0011"+
		"\u0007\u010c\u010d\n\u0005\u0000\u0000\u010d\u010e\u0007\u0006\u0000\u0000"+
		"\u010e\u0119\u0003\"\u0011\u0006\u010f\u0110\n\u0004\u0000\u0000\u0110"+
		"\u0111\u0005\u001f\u0000\u0000\u0111\u0119\u0003\"\u0011\u0005\u0112\u0113"+
		"\n\u0003\u0000\u0000\u0113\u0114\u0005 \u0000\u0000\u0114\u0119\u0003"+
		"\"\u0011\u0004\u0115\u0116\n\u0001\u0000\u0000\u0116\u0117\u0005*\u0000"+
		"\u0000\u0117\u0119\u00051\u0000\u0000\u0118\u0103\u0001\u0000\u0000\u0000"+
		"\u0118\u0106\u0001\u0000\u0000\u0000\u0118\u0109\u0001\u0000\u0000\u0000"+
		"\u0118\u010c\u0001\u0000\u0000\u0000\u0118\u010f\u0001\u0000\u0000\u0000"+
		"\u0118\u0112\u0001\u0000\u0000\u0000\u0118\u0115\u0001\u0000\u0000\u0000"+
		"\u0119\u011c\u0001\u0000\u0000\u0000\u011a\u0118\u0001\u0000\u0000\u0000"+
		"\u011a\u011b\u0001\u0000\u0000\u0000\u011b#\u0001\u0000\u0000\u0000\u011c"+
		"\u011a\u0001\u0000\u0000\u0000\u011d\u0122\u0003&\u0013\u0000\u011e\u011f"+
		"\u0005\u0005\u0000\u0000\u011f\u0121\u0003&\u0013\u0000\u0120\u011e\u0001"+
		"\u0000\u0000\u0000\u0121\u0124\u0001\u0000\u0000\u0000\u0122\u0120\u0001"+
		"\u0000\u0000\u0000\u0122\u0123\u0001\u0000\u0000\u0000\u0123\u0126\u0001"+
		"\u0000\u0000\u0000\u0124\u0122\u0001\u0000\u0000\u0000\u0125\u0127\u0003"+
		"(\u0014\u0000\u0126\u0125\u0001\u0000\u0000\u0000\u0126\u0127\u0001\u0000"+
		"\u0000\u0000\u0127%\u0001\u0000\u0000\u0000\u0128\u0129\u0005\u0004\u0000"+
		"\u0000\u0129\u012a\u0003\"\u0011\u0000\u012a\u012e\u0005$\u0000\u0000"+
		"\u012b\u012d\u0003\u0002\u0001\u0000\u012c\u012b\u0001\u0000\u0000\u0000"+
		"\u012d\u0130\u0001\u0000\u0000\u0000\u012e\u012c\u0001\u0000\u0000\u0000"+
		"\u012e\u012f\u0001\u0000\u0000\u0000\u012f\u0131\u0001\u0000\u0000\u0000"+
		"\u0130\u012e\u0001\u0000\u0000\u0000\u0131\u0132\u0005%\u0000\u0000\u0132"+
		"\'\u0001\u0000\u0000\u0000\u0133\u0134\u0005\u0005\u0000\u0000\u0134\u0138"+
		"\u0005$\u0000\u0000\u0135\u0137\u0003\u0002\u0001\u0000\u0136\u0135\u0001"+
		"\u0000\u0000\u0000\u0137\u013a\u0001\u0000\u0000\u0000\u0138\u0136\u0001"+
		"\u0000\u0000\u0000\u0138\u0139\u0001\u0000\u0000\u0000\u0139\u013b\u0001"+
		"\u0000\u0000\u0000\u013a\u0138\u0001\u0000\u0000\u0000\u013b\u013c\u0005"+
		"%\u0000\u0000\u013c)\u0001\u0000\u0000\u0000\u013d\u013e\u0005\u0006\u0000"+
		"\u0000\u013e\u013f\u0003\"\u0011\u0000\u013f\u0143\u0005$\u0000\u0000"+
		"\u0140\u0142\u0003,\u0016\u0000\u0141\u0140\u0001\u0000\u0000\u0000\u0142"+
		"\u0145\u0001\u0000\u0000\u0000\u0143\u0141\u0001\u0000\u0000\u0000\u0143"+
		"\u0144\u0001\u0000\u0000\u0000\u0144\u0147\u0001\u0000\u0000\u0000\u0145"+
		"\u0143\u0001\u0000\u0000\u0000\u0146\u0148\u0003.\u0017\u0000\u0147\u0146"+
		"\u0001\u0000\u0000\u0000\u0147\u0148\u0001\u0000\u0000\u0000\u0148\u0149"+
		"\u0001\u0000\u0000\u0000\u0149\u014a\u0005%\u0000\u0000\u014a+\u0001\u0000"+
		"\u0000\u0000\u014b\u014c\u0005\u0007\u0000\u0000\u014c\u014d\u0003\"\u0011"+
		"\u0000\u014d\u0151\u0005)\u0000\u0000\u014e\u0150\u0003\u0002\u0001\u0000"+
		"\u014f\u014e\u0001\u0000\u0000\u0000\u0150\u0153\u0001\u0000\u0000\u0000"+
		"\u0151\u014f\u0001\u0000\u0000\u0000\u0151\u0152\u0001\u0000\u0000\u0000"+
		"\u0152-\u0001\u0000\u0000\u0000\u0153\u0151\u0001\u0000\u0000\u0000\u0154"+
		"\u0155\u0005\b\u0000\u0000\u0155\u0159\u0005)\u0000\u0000\u0156\u0158"+
		"\u0003\u0002\u0001\u0000\u0157\u0156\u0001\u0000\u0000\u0000\u0158\u015b"+
		"\u0001\u0000\u0000\u0000\u0159\u0157\u0001\u0000\u0000\u0000\u0159\u015a"+
		"\u0001\u0000\u0000\u0000\u015a/\u0001\u0000\u0000\u0000\u015b\u0159\u0001"+
		"\u0000\u0000\u0000\u015c\u015d\u0005\n\u0000\u0000\u015d\u015e\u0003\""+
		"\u0011\u0000\u015e\u0162\u0005$\u0000\u0000\u015f\u0161\u0003\u0002\u0001"+
		"\u0000\u0160\u015f\u0001\u0000\u0000\u0000\u0161\u0164\u0001\u0000\u0000"+
		"\u0000\u0162\u0160\u0001\u0000\u0000\u0000\u0162\u0163\u0001\u0000\u0000"+
		"\u0000\u0163\u0165\u0001\u0000\u0000\u0000\u0164\u0162\u0001\u0000\u0000"+
		"\u0000\u0165\u0166\u0005%\u0000\u0000\u01661\u0001\u0000\u0000\u0000\u0167"+
		"\u0168\u0005\t\u0000\u0000\u0168\u0169\u0003\"\u0011\u0000\u0169\u016d"+
		"\u0005$\u0000\u0000\u016a\u016c\u0003\u0002\u0001\u0000\u016b\u016a\u0001"+
		"\u0000\u0000\u0000\u016c\u016f\u0001\u0000\u0000\u0000\u016d\u016b\u0001"+
		"\u0000\u0000\u0000\u016d\u016e\u0001\u0000\u0000\u0000\u016e\u0170\u0001"+
		"\u0000\u0000\u0000\u016f\u016d\u0001\u0000\u0000\u0000\u0170\u0171\u0005"+
		"%\u0000\u0000\u0171\u0191\u0001\u0000\u0000\u0000\u0172\u0173\u0005\t"+
		"\u0000\u0000\u0173\u0174\u0003\u001a\r\u0000\u0174\u0175\u0005(\u0000"+
		"\u0000\u0175\u0176\u0003\"\u0011\u0000\u0176\u0177\u0005(\u0000\u0000"+
		"\u0177\u0178\u0003\"\u0011\u0000\u0178\u017c\u0005$\u0000\u0000\u0179"+
		"\u017b\u0003\u0002\u0001\u0000\u017a\u0179\u0001\u0000\u0000\u0000\u017b"+
		"\u017e\u0001\u0000\u0000\u0000\u017c\u017a\u0001\u0000\u0000\u0000\u017c"+
		"\u017d\u0001\u0000\u0000\u0000\u017d\u017f\u0001\u0000\u0000\u0000\u017e"+
		"\u017c\u0001\u0000\u0000\u0000\u017f\u0180\u0005%\u0000\u0000\u0180\u0191"+
		"\u0001\u0000\u0000\u0000\u0181\u0182\u0005\t\u0000\u0000\u0182\u0183\u0005"+
		"1\u0000\u0000\u0183\u0184\u0005+\u0000\u0000\u0184\u0185\u00051\u0000"+
		"\u0000\u0185\u0186\u0005\u000b\u0000\u0000\u0186\u0187\u0003\"\u0011\u0000"+
		"\u0187\u018b\u0005$\u0000\u0000\u0188\u018a\u0003\u0002\u0001\u0000\u0189"+
		"\u0188\u0001\u0000\u0000\u0000\u018a\u018d\u0001\u0000\u0000\u0000\u018b"+
		"\u0189\u0001\u0000\u0000\u0000\u018b\u018c\u0001\u0000\u0000\u0000\u018c"+
		"\u018e\u0001\u0000\u0000\u0000\u018d\u018b\u0001\u0000\u0000\u0000\u018e"+
		"\u018f\u0005%\u0000\u0000\u018f\u0191\u0001\u0000\u0000\u0000\u0190\u0167"+
		"\u0001\u0000\u0000\u0000\u0190\u0172\u0001\u0000\u0000\u0000\u0190\u0181"+
		"\u0001\u0000\u0000\u0000\u01913\u0001\u0000\u0000\u0000\u0192\u0193\u0003"+
		"\"\u0011\u0000\u0193\u0194\u0005*\u0000\u0000\u0194\u0195\u0005*\u0000"+
		"\u0000\u0195\u0196\u0005*\u0000\u0000\u0196\u0197\u0003\"\u0011\u0000"+
		"\u01975\u0001\u0000\u0000\u0000\u0198\u019a\u0005\u000e\u0000\u0000\u0199"+
		"\u019b\u0003\"\u0011\u0000\u019a\u0199\u0001\u0000\u0000\u0000\u019a\u019b"+
		"\u0001\u0000\u0000\u0000\u019b\u019f\u0001\u0000\u0000\u0000\u019c\u019f"+
		"\u0005\f\u0000\u0000\u019d\u019f\u0005\r\u0000\u0000\u019e\u0198\u0001"+
		"\u0000\u0000\u0000\u019e\u019c\u0001\u0000\u0000\u0000\u019e\u019d\u0001"+
		"\u0000\u0000\u0000\u019f7\u0001\u0000\u0000\u0000\u01a0\u01a1\u0003\u001c"+
		"\u000e\u0000\u01a1\u01a3\u0005\"\u0000\u0000\u01a2\u01a4\u0003<\u001e"+
		"\u0000\u01a3\u01a2\u0001\u0000\u0000\u0000\u01a3\u01a4\u0001\u0000\u0000"+
		"\u0000\u01a4\u01a5\u0001\u0000\u0000\u0000\u01a5\u01a6\u0005#\u0000\u0000"+
		"\u01a69\u0001\u0000\u0000\u0000\u01a7\u01ab\u0005$\u0000\u0000\u01a8\u01aa"+
		"\u0003\u0002\u0001\u0000\u01a9\u01a8\u0001\u0000\u0000\u0000\u01aa\u01ad"+
		"\u0001\u0000\u0000\u0000\u01ab\u01a9\u0001\u0000\u0000\u0000\u01ab\u01ac"+
		"\u0001\u0000\u0000\u0000\u01ac\u01ae\u0001\u0000\u0000\u0000\u01ad\u01ab"+
		"\u0001\u0000\u0000\u0000\u01ae\u01af\u0005%\u0000\u0000\u01af;\u0001\u0000"+
		"\u0000\u0000\u01b0\u01b5\u0003>\u001f\u0000\u01b1\u01b2\u0005+\u0000\u0000"+
		"\u01b2\u01b4\u0003>\u001f\u0000\u01b3\u01b1\u0001\u0000\u0000\u0000\u01b4"+
		"\u01b7\u0001\u0000\u0000\u0000\u01b5\u01b3\u0001\u0000\u0000\u0000\u01b5"+
		"\u01b6\u0001\u0000\u0000\u0000\u01b6=\u0001\u0000\u0000\u0000\u01b7\u01b5"+
		"\u0001\u0000\u0000\u0000\u01b8\u01ba\u00051\u0000\u0000\u01b9\u01b8\u0001"+
		"\u0000\u0000\u0000\u01b9\u01ba\u0001\u0000\u0000\u0000\u01ba\u01bd\u0001"+
		"\u0000\u0000\u0000\u01bb\u01be\u0003\u001c\u000e\u0000\u01bc\u01be\u0003"+
		"\"\u0011\u0000\u01bd\u01bb\u0001\u0000\u0000\u0000\u01bd\u01bc\u0001\u0000"+
		"\u0000\u0000\u01be?\u0001\u0000\u0000\u0000\u01bf\u01c0\u0005\u0002\u0000"+
		"\u0000\u01c0\u01c1\u00051\u0000\u0000\u01c1\u01c3\u0005\"\u0000\u0000"+
		"\u01c2\u01c4\u0003B!\u0000\u01c3\u01c2\u0001\u0000\u0000\u0000\u01c3\u01c4"+
		"\u0001\u0000\u0000\u0000\u01c4\u01c5\u0001\u0000\u0000\u0000\u01c5\u01c7"+
		"\u0005#\u0000\u0000\u01c6\u01c8\u0003\u0018\f\u0000\u01c7\u01c6\u0001"+
		"\u0000\u0000\u0000\u01c7\u01c8\u0001\u0000\u0000\u0000\u01c8\u01c9\u0001"+
		"\u0000\u0000\u0000\u01c9\u01cd\u0005$\u0000\u0000\u01ca\u01cc\u0003\u0002"+
		"\u0001\u0000\u01cb\u01ca\u0001\u0000\u0000\u0000\u01cc\u01cf\u0001\u0000"+
		"\u0000\u0000\u01cd\u01cb\u0001\u0000\u0000\u0000\u01cd\u01ce\u0001\u0000"+
		"\u0000\u0000\u01ce\u01d0\u0001\u0000\u0000\u0000\u01cf\u01cd\u0001\u0000"+
		"\u0000\u0000\u01d0\u01d1\u0005%\u0000\u0000\u01d1A\u0001\u0000\u0000\u0000"+
		"\u01d2\u01d7\u0003D\"\u0000\u01d3\u01d4\u0005+\u0000\u0000\u01d4\u01d6"+
		"\u0003D\"\u0000\u01d5\u01d3\u0001\u0000\u0000\u0000\u01d6\u01d9\u0001"+
		"\u0000\u0000\u0000\u01d7\u01d5\u0001\u0000\u0000\u0000\u01d7\u01d8\u0001"+
		"\u0000\u0000\u0000\u01d8C\u0001\u0000\u0000\u0000\u01d9\u01d7\u0001\u0000"+
		"\u0000\u0000\u01da\u01db\u00051\u0000\u0000\u01db\u01dc\u0003\u0018\f"+
		"\u0000\u01dcE\u0001\u0000\u0000\u0000\u01dd\u01de\u0005\u0003\u0000\u0000"+
		"\u01de\u01df\u00051\u0000\u0000\u01df\u01e1\u0005$\u0000\u0000\u01e0\u01e2"+
		"\u0003H$\u0000\u01e1\u01e0\u0001\u0000\u0000\u0000\u01e2\u01e3\u0001\u0000"+
		"\u0000\u0000\u01e3\u01e1\u0001\u0000\u0000\u0000\u01e3\u01e4\u0001\u0000"+
		"\u0000\u0000\u01e4\u01e5\u0001\u0000\u0000\u0000\u01e5\u01e6\u0005%\u0000"+
		"\u0000\u01e6G\u0001\u0000\u0000\u0000\u01e7\u01e8\u0003\u0018\f\u0000"+
		"\u01e8\u01e9\u00051\u0000\u0000\u01e9\u01ea\u0005(\u0000\u0000\u01eaI"+
		"\u0001\u0000\u0000\u0000\u01eb\u01f0\u0003L&\u0000\u01ec\u01ed\u0005+"+
		"\u0000\u0000\u01ed\u01ef\u0003L&\u0000\u01ee\u01ec\u0001\u0000\u0000\u0000"+
		"\u01ef\u01f2\u0001\u0000\u0000\u0000\u01f0\u01ee\u0001\u0000\u0000\u0000"+
		"\u01f0\u01f1\u0001\u0000\u0000\u0000\u01f1\u01f4\u0001\u0000\u0000\u0000"+
		"\u01f2\u01f0\u0001\u0000\u0000\u0000\u01f3\u01f5\u0005+\u0000\u0000\u01f4"+
		"\u01f3\u0001\u0000\u0000\u0000\u01f4\u01f5\u0001\u0000\u0000\u0000\u01f5"+
		"K\u0001\u0000\u0000\u0000\u01f6\u01f7\u00051\u0000\u0000\u01f7\u01f8\u0005"+
		")\u0000\u0000\u01f8\u01f9\u0003\"\u0011\u0000\u01f9M\u0001\u0000\u0000"+
		"\u0000-QUc\u0083\u008d\u0090\u009b\u00a7\u00bd\u00c6\u00d4\u00db\u00e3"+
		"\u00e9\u00fe\u0101\u0118\u011a\u0122\u0126\u012e\u0138\u0143\u0147\u0151"+
		"\u0159\u0162\u016d\u017c\u018b\u0190\u019a\u019e\u01a3\u01ab\u01b5\u01b9"+
		"\u01bd\u01c3\u01c7\u01cd\u01d7\u01e3\u01f0\u01f4";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}