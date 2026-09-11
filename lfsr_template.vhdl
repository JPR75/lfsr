-- --------------------------------------------------------------------------------
-- File        : lfsr_template.vhdl
-- Project     : LFSR
-- Description : VHDL template to generate a frequency divider based on LFSR
--               {!size}  will be rplaced by size of the polynomial - 1
--               {!count} will be rplaced by size of the taget count
--               {!polynomial} will be replaced by the polynomial
--
--               example: division by 100000 => P = X^17 + X^14 ; LFSR target count = 0x546B
--                                           => {!size} = 16
--                                           => {!count} = "00101010001101011"
--                                           => {!polynomial"} = s_lfsr(16) xnor s_lfsr(13)
--
-- Note: Don't forget to replace this header with your own!
--
-- --------------------------------------------------------------------------------
-- Author      : JPR75 (https://github.com/JPR75/lfsr.git)
-- --------------------------------------------------------------------------------
-- Copyright (C) 2020 - 2026 JPR75
--
-- This program is free software: you can redistribute it and/or modify
-- it under the terms of the GNU Lesser General Public License as published by
-- the Free Software Foundation, either version 3 of the License, or
-- (at your option) any later version.
--
-- This program is distributed in the hope that it will be useful,
-- but WITHOUT ANY WARRANTY; without even the implied warranty of
-- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
-- GNU Lesser General Public License for more details.
--
-- You should have received a copy of the GNU Lesser General Public License
-- along with this program.  If not, see <http://www.gnu.org/licenses/>
-- --------------------------------------------------------------------------------

library IEEE;
use IEEE.STD_LOGIC_1164.ALL;

library work;

entity lfsr is
    port (
        clk_in    : in  std_logic;
        reset     : in  std_logic;
        pulse_out : out std_logic
    );
end entity lfsr;

architecture rtl_lfsr of lfsr is
    signal s_lfsr : std_logic_vector({!size} downto 0) := (others => '0');
    signal s_clk  : std_logic := '0';

    begin

    lfsr: process(clk_in)
        variable s_feedback : std_logic;
    begin
        if rising_edge(clk_in) then
            if reset = '1' then
                s_lfsr <= (others => '0');
                s_clk  <= '1';
            elsif s_lfsr = "{!count}" then
                s_lfsr <= (OTHERS => '0');
                s_clk  <= '1';
            else
                s_feedback := {!polynomial};
                s_lfsr <= s_lfsr(({!size} - 1) downto 0) & s_feedback;
                s_clk  <= '0';
            end if;
        end if;
    end process;

  pulse_out  <= s_clk;

end architecture rtl_lfsr;
